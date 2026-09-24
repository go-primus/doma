package svc

import (
	"context"
	"log/slog"
	"reflect"
)

var sm ServiceManager

func init() {
	sm = newServiceRegistry()
}

func GetInstance() ServiceManager {
	return sm
}

func GetSystemService[T Service](name string) T {
	svc := GetInstance().GetSystemService(name)
	return svc.(T)
}

func GetSystemServiceState(name string) (ServiceState, error) {
	return GetInstance().GetServiceState(name)
}

type ServiceManager interface {
	// GetService()
	// CheckService()
	// AddService()
	// ListServices()
	// GetServiceDebugInfo()

	RegisterService(srv Service)
	Init(ctx context.Context) error
	Config(ctx context.Context) error
	Start(ctx context.Context) error
	Stop(ctx context.Context) error

	GetSystemService(name string) Service
	GetServiceState(name string) (ServiceState, error)
	IsSystemReady() bool
	DumpServices()
}

type ServiceRegistry struct {
	services map[string]Service
	isready  bool
}

// IsSystemReady implements [ServiceManager].
func (s *ServiceRegistry) IsSystemReady() bool {
	return s.isready
}

func newServiceRegistry() ServiceManager {
	return &ServiceRegistry{
		services: make(map[string]Service),
	}
}

func (s *ServiceRegistry) RegisterService(srv Service) {
	s.services[srv.Name()] = srv
}

func (s *ServiceRegistry) Init(ctx context.Context) error {
	slog.Info("init services...")
	for _, svc := range s.services {
		err := svc.Init(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *ServiceRegistry) Config(ctx context.Context) error {
	slog.Info("config services...")
	for _, svc := range s.services {
		err := svc.Config(ctx, func() error { return nil })
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *ServiceRegistry) Start(ctx context.Context) error {
	slog.Info("start services...")
	for _, svc := range s.services {
		err := svc.Start(ctx)
		if err != nil {
			return err
		}
	}
	s.isready = true

	for _, svc := range s.services {
		if h, ok := svc.(LifecycleSystemReadyHook); ok {
			if err := h.OnSystemReady(ctx); err != nil {
				return err
			}
		}
	}

	return nil
}

// Stop implements [ServiceManager].
func (s *ServiceRegistry) Stop(ctx context.Context) error {
	slog.Info("stop services...")
	for _, svc := range s.services {
		err := svc.Stop(ctx)
		if err != nil {
			slog.Warn("stop service fail", "service", svc.Name(), "err", err)
		}
	}
	s.isready = false
	return nil
}

///

func (s *ServiceRegistry) GetSystemService(name string) Service {
	return s.services[name]
}

func (s *ServiceRegistry) GetServiceState(name string) (ServiceState, error) {
	return s.GetSystemService(name).Status()
}

func (s *ServiceRegistry) DumpService(name string) {
	srv := s.GetSystemService(name)
	status, _ := srv.Status()
	slog.Info("service info:", "srvid", srv.ID(), "name", srv.Name(), "status", status, "type", reflect.TypeOf(srv))

}

func (s *ServiceRegistry) DumpServices() {
	slog.Info("---dump services------------------------")
	for _, svc := range s.services {
		s.DumpService(svc.Name())
	}
	slog.Info("---------------------------------------")
}
