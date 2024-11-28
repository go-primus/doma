package svc

import (
	"context"

	"github.com/primus/primus/core/common/command"
	"github.com/primus/primus/core/common/event"
	"github.com/primus/primus/core/common/task"
)

//go:generate stringer -type=ServiceState
type ServiceState int32

// Service State: Idle -> Init -> Config -> Ready
const (
	ServiceState_Idle    ServiceState = iota // 空闲
	ServiceState_Init                        // 初始化
	ServiceState_Config                      // 配置中
	ServiceState_Ready                       // 就绪
	ServiceState_Migrate                     // 迁移
	ServiceState_Running                     // 运行，提供服务
	ServiceState_Pause                       // 暂停
	ServiceState_Destroy                     // 销毁
)

type ServiceHandler interface {
	OnStart()
	OnSystemReady()
	OnStop()
	OnDestory()
}

type Service interface {
	ID() string
	Name() string
	Version() string

	Init(context.Context) error
	Config(context.Context, func() error) error //load config
	Start(context.Context) error
	// Pause(context.Context) error
	Stop(context.Context) error    //
	Status() (ServiceState, error) // 获取服务状态
	Dump(context.Context) ([]byte, error)

	// command
	HandleCommand(context.Context, command.Command) error
	SubscribeCommand(fn CommandHandler) error

	// Event
	Publish(topic string, data event.Event) error
	Subscribe(topic string, fn EventHandler) error

	ServiceTask
}

type ProcessWithKey interface {
	KvStore
	ProcessWithKey(ctx context.Context, key string, fn func(ctx context.Context) error) error
}

type ServiceTask interface {
	StartTask(target string, task task.Task) error
	SubscribeTask(TaskHandler) error
}

type CommandHandler func(cmd command.Command) (any, error)
type EventHandler func(ctx context.Context, topic string, event event.Event) error
type TaskHandler func(ctx context.Context, task task.Task) (any, error)

type ServiceHandle interface {
	OnCommand(cmd command.Command) (any, error)
	OnEvent(topic string, event event.Event) error
}
