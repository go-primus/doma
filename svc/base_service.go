package svc

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"

	"github.com/go-primus/doma/core/common/command"
	"github.com/go-primus/doma/core/common/event"
	"github.com/go-primus/doma/pkg/eventbus"
	"github.com/google/uuid"
)

type BaseService struct {
	id        string
	name      string
	version   string
	state     int32
	bus       eventbus.EventBus
	lifecycle ServiceLifecycle

	handle ServiceHandle

	// tasks map[string]struct{}
	tasks sync.Map
}

func NewBaseService(name string) *BaseService {
	return newBaseService(name, nil)
}

func NewBaseServiceWithBus(name string, bus eventbus.EventBus) *BaseService {
	return newBaseService(name, bus)
}

// BindLifecycle 绑定业务服务生命周期回调：依据业务服务实现的可选接口
// （LifecycleInitHook/LifecycleConfigHook/LifecycleStartHook/LifecycleStopHook）
// 构造 ServiceLifecycle 默认实现并绑定，未实现的阶段为无操作。
func (s *BaseService) BindLifecycle(provider Service) {
	s.lifecycle = newServiceLifecycle(provider)
}

func newBaseService(name string, bus eventbus.EventBus) *BaseService {

	//

	svc := &BaseService{
		id:   uuid.NewString(),
		name: name,
		bus:  bus,
	}
	// svc.Init(context.Background())
	return svc
}

func (s *BaseService) ID() string {
	return s.id
}

func (s *BaseService) Name() string {
	return s.name
}

func (s *BaseService) Version() string {
	return s.version
}

func (s *BaseService) getState() ServiceState {
	return ServiceState(atomic.LoadInt32(&s.state))
}

func (s *BaseService) updateState(state ServiceState) {
	atomic.StoreInt32(&s.state, int32(state))

}

//////////////

func (s *BaseService) Init(ctx context.Context) error {
	if s.getState() != ServiceState_Idle {
		return fmt.Errorf("service has inited")
	}
	if s.lifecycle != nil {
		if err := s.lifecycle.OnInit(ctx); err != nil {
			return err
		}
	}
	s.updateState(ServiceState_Init)
	return nil
}

func (s *BaseService) Config(ctx context.Context, fn func() error) error {
	if s.getState() != ServiceState_Init {
		if s.getState() == ServiceState_Idle {
			return fmt.Errorf("should init service before config")
		}
		return fmt.Errorf("service has configed")
	}
	if s.lifecycle != nil {
		if err := s.lifecycle.OnConfig(ctx); err != nil {
			return err
		}
	}

	// s.state = ServiceState_Config
	s.updateState(ServiceState_Config)

	if err := fn(); err != nil {
		return fmt.Errorf("config failed ,%w", err) //errors.Wrap(err, "config failed")
	}
	// s.state = ServiceState_Ready
	s.updateState(ServiceState_Ready)
	return nil
}

func (s *BaseService) Start(ctx context.Context) error {
	if s.getState() != ServiceState_Ready {
		return fmt.Errorf("service not ready")
	}
	slog.Info("service start :", "service", s.name)
	if s.lifecycle != nil {
		if err := s.lifecycle.OnStart(ctx); err != nil {
			return err
		}
	}
	// s.state = ServiceState_Running
	s.updateState(ServiceState_Running)

	return nil
}

// Stop implements [Service].
func (s *BaseService) Stop(ctx context.Context) error {
	if s.getState() != ServiceState_Running {
		return fmt.Errorf("service not running")
	}
	if s.lifecycle != nil {
		if err := s.lifecycle.OnStop(ctx); err != nil {
			slog.Warn("stop lifecycle failed", "service", s.name, "err", err)
		}
	}
	// s.state = ServiceState_Destroy
	s.updateState(ServiceState_Destroy)
	return nil
}

func (s *BaseService) Status() (ServiceState, error) {
	return s.getState(), nil
}

func (s *BaseService) Dump(context.Context) ([]byte, error) {

	return nil, fmt.Errorf("service %s dump not implemented", s.name)
}

// /////////////
func (s *BaseService) Publish(topic string, data event.Event) error {
	if s.bus == nil {
		slog.Warn("publish err: bus is nil", "service", s.name)
		return nil
	}
	slog.Info("publish topic ", "service", s.name, "topic", topic)
	return s.bus.Publish(topic, data)
}

func (s *BaseService) Subscribe(topic string, fn EventHandler) error {
	if s.bus == nil {
		return nil
	}

	s.bus.Subscribe(topic, s.handleMsg(fn))
	return nil
}

func (s *BaseService) handleMsg(fn EventHandler) eventbus.EventHandler {
	return func(msg *eventbus.Msg) {

		if msg == nil {
			return
		}

		event := event.Event{}
		err := json.Unmarshal(msg.Data, &event)
		if err != nil {
			slog.Error("invalid event,", "err", err, "data", string(msg.Data))
			return
		}

		err = fn(context.Background(), msg.Subject, event)
		if err != nil {
			slog.Error("handle message ", "service", s.name, "subject", msg.Subject, "err", err)
			return
		}
	}
}

// ///
func (s *BaseService) SubscribeCommand(fn CommandHandler) error {
	// s.Subscribe("command.file.>")
	if s.bus == nil {
		return nil
	}
	if fn == nil {
		return nil
	}

	s.bus.Subscribe(fmt.Sprintf("command.%s", s.name), s.handleCommand(fn))
	return nil
}

func (s *BaseService) handleCommand(fn CommandHandler) eventbus.EventHandler {
	return func(msg *eventbus.Msg) {
		cmd := command.Command{}
		err := json.Unmarshal(msg.Data, &cmd)
		if err != nil {
			slog.Error("handle command, unmarshal err:", "service", s.name, "err", err)
			return
		}

		slog.Info("receive command bus msg: ", "service", s.name, "cmd", cmd.CommandType)
		res, err := fn(cmd)
		if err != nil {
			slog.Error("handle command ", "service", s.name, "cmd", cmd.CommandType, " err", err)
			s.bus.Publish(msg.Reply, err.Error())
			return
		}
		bb, _ := json.Marshal(res)
		s.bus.Publish(msg.Reply, string(bb))
	}
}

func (s *BaseService) HandleCommand(ctx context.Context, cmd command.Command) error {
	//
	// a,err :=  s.store.Load(ctx, aggtype, cmd.ID())
	// a.HandleCommand(ctx,cmd)
	// a.store.Save(ctx,a)
	return nil
}
