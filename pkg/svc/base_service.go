package svc

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/primus/primus/core/common/command"
	"github.com/primus/primus/core/common/event"
	"github.com/primus/primus/pkg/eventbus"
	"github.com/primus/primus/pkg/logger"
	"github.com/sirupsen/logrus"
)

type BaseService struct {
	id      string
	name    string
	version string
	state   int32
	bus     eventbus.EventBus

	handle ServiceHandle

	// tasks map[string]struct{}
	tasks sync.Map
}

func NewBaseService(name string) Service {
	return newBaseService(name, nil)
}

func NewBaseServiceWithBus(name string, bus eventbus.EventBus) Service {
	return newBaseService(name, bus)
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

func (s *BaseService) Init(context.Context) error {
	if s.getState() != ServiceState_Idle {
		return fmt.Errorf("service has inited")
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

	// s.state = ServiceState_Config
	s.updateState(ServiceState_Config)

	if err := fn(); err != nil {
		return fmt.Errorf("config failed ,%w", err) //errors.Wrap(err, "config failed")
	}
	// s.state = ServiceState_Ready
	s.updateState(ServiceState_Ready)
	return nil
}

func (s *BaseService) Start(context.Context) error {
	if s.getState() != ServiceState_Ready {
		return fmt.Errorf("service not ready")
	}
	logrus.Info("service start :", s.name)
	// s.state = ServiceState_Running
	s.updateState(ServiceState_Running)

	return nil
}

func (s *BaseService) Stop(context.Context) error {
	if s.getState() != ServiceState_Running {
		return fmt.Errorf("service not running")
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
		logger.L().Warn("----service:", s.name, ", publish err: bus is nil")
		return nil
	}
	logger.L().Info("----service:", s.name, ", publish ", topic)
	return s.bus.Publish(topic, data)
}

func (s *BaseService) Subscribe(topic string, fn EventHandler) error {
	if s.bus == nil {
		return nil
	}

	s.bus.Subscribe(topic, s.handleMsg(fn))
	return nil
}

func (s *BaseService) handleMsg(fn EventHandler) func(msg *nats.Msg) {
	return func(msg *nats.Msg) {

		if msg == nil {
			return
		}

		event := event.Event{}
		err := json.Unmarshal(msg.Data, &event)
		if err != nil {
			logger.L().Error("invalid event,", err, ", ", string(msg.Data))
			return
		}

		err = fn(context.Background(), msg.Subject, event)
		if err != nil {
			logger.L().Error("----service ", s.name, "---- handle message ", msg.Subject, ", err:", err)
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

func (s *BaseService) handleCommand(fn CommandHandler) func(msg *nats.Msg) {
	return func(msg *nats.Msg) {
		cmd := command.Command{}
		err := json.Unmarshal(msg.Data, &cmd)
		if err != nil {
			logger.L().Error("----service:", s.name, ", handle command, unmarshal err:", err)
			return
		}

		logger.L().Info("----service:", s.name, ", receive command bus msg: ", cmd.CommandType)
		res, err := fn(cmd)
		if err != nil {
			logger.L().Error("----service:", s.name, ", handle command ", cmd.CommandType, " err:", err)
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
