package svc

import "context"

// 可选生命周期接口：业务服务按需实现其一或多个，经 BaseService.BindLifecycle 绑定后，
// 由 BaseService 在状态机推进前调用对应回调，回调成功后才更新状态，未实现的阶段自动跳过。
type LifecycleInitHook interface {
	OnInit(ctx context.Context) error
}

type LifecycleConfigHook interface {
	OnConfig(ctx context.Context) error
}

type LifecycleStartHook interface {
	OnStart(ctx context.Context) error
}

type LifecycleStopHook interface {
	OnStop(ctx context.Context) error
}

// LifecycleSystemReadyHook 系统就绪通知：由 ServiceManager 在所有服务 Start 完成、
// isready 置位后的第二遍统一调用，未实现的服务自动跳过。
// 不入 BindLifecycle/serviceLifecycle：OnSystemReady 曾为 Service 接口必选方法，
// 经嵌入提升后断言会绑定到基类自身方法导致递归。
type LifecycleSystemReadyHook interface {
	OnSystemReady(ctx context.Context) error
}

// ServiceLifecycle 完整服务生命周期钩子接口：由 BindLifecycle 依据业务服务实现的
// 可选接口自动构造默认实现 serviceLifecycle。
type ServiceLifecycle interface {
	OnInit(ctx context.Context) error
	OnConfig(ctx context.Context) error
	OnStart(ctx context.Context) error
	OnStop(ctx context.Context) error
}

// serviceLifecycle ServiceLifecycle 的默认实现：BindLifecycle 时按业务服务实现的
// 可选接口填充对应函数字段，未实现的阶段保持 nil（方法直接返回 nil）。
type serviceLifecycle struct {
	onInit   func(ctx context.Context) error
	onConfig func(ctx context.Context) error
	onStart  func(ctx context.Context) error
	onStop   func(ctx context.Context) error
}

var _ ServiceLifecycle = (*serviceLifecycle)(nil)

// newServiceLifecycle 依据 provider 实现的可选接口构造 ServiceLifecycle 默认实现。
func newServiceLifecycle(provider Service) ServiceLifecycle {
	lc := &serviceLifecycle{}
	if h, ok := provider.(LifecycleInitHook); ok {
		lc.onInit = h.OnInit
	}
	if h, ok := provider.(LifecycleConfigHook); ok {
		lc.onConfig = h.OnConfig
	}
	if h, ok := provider.(LifecycleStartHook); ok {
		lc.onStart = h.OnStart
	}
	if h, ok := provider.(LifecycleStopHook); ok {
		lc.onStop = h.OnStop
	}
	return lc
}

func (l *serviceLifecycle) OnInit(ctx context.Context) error {
	if l.onInit == nil {
		return nil
	}
	return l.onInit(ctx)
}

func (l *serviceLifecycle) OnConfig(ctx context.Context) error {
	if l.onConfig == nil {
		return nil
	}
	return l.onConfig(ctx)
}

func (l *serviceLifecycle) OnStart(ctx context.Context) error {
	if l.onStart == nil {
		return nil
	}
	return l.onStart(ctx)
}

func (l *serviceLifecycle) OnStop(ctx context.Context) error {
	if l.onStop == nil {
		return nil
	}
	return l.onStop(ctx)
}
