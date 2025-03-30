package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-primus/doma/runtime/aggregatestore"
	"github.com/go-primus/doma/runtime/core"
	"github.com/go-primus/doma/runtime/eventstore"
	memoryEventStore "github.com/go-primus/doma/runtime/eventstore/memory"
	"github.com/go-primus/doma/runtime/registry"
	"github.com/google/uuid"
)

type DomaDemoService interface {
	CreateDemo(uuid.UUID, string) error
}

type demoDemoServiceImpl struct {
	t     core.AggregateType
	store aggregatestore.AggregateStore
}

func (d *demoDemoServiceImpl) handleCommand(ctx context.Context, cmd core.Command) error {
	//
	a, err := d.store.Load(ctx, d.t, cmd.AggregateID())
	if err != nil {
		fmt.Println("handle command load err:", err)
		return err
	}
	if err = a.HandleCommand(ctx, cmd); err != nil {
		return err
	}

	return d.store.Save(ctx, a)
}

// CreateDemo implements DomaDemoService.
func (d *demoDemoServiceImpl) CreateDemo(uuidx uuid.UUID, name string) error {

	return d.handleCommand(context.Background(), &CreateDemo{
		ID:   uuidx,
		Name: name,
	})
}

func NewDemoService(estore eventstore.EventStore) DomaDemoService {
	var store aggregatestore.AggregateStore
	store, _ = aggregatestore.NewAggregateStore(estore)
	return &demoDemoServiceImpl{
		t:     TestAggregateRegisterType,
		store: store,
	}
}

func TestDoma(t *testing.T) {

	// 注册聚合

	registry.RegisterAggregate(func(u uuid.UUID) core.Aggregate {
		return NewInvitationAggregate(u)
	})

	// 绑定关系

	eventStore, _ := memoryEventStore.NewEventStore(
	// memoryEventStore.WithEventHandler(eventBus), // Add the event bus as a handler after save.
	)

	// 发送命令
	id := uuid.New()
	s := NewDemoService(eventStore)
	err := s.CreateDemo(id, "name")
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println("------------")

	err = s.CreateDemo(id, "name")
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println("------------")
	err = s.CreateDemo(id, "name")
	if err != nil {
		t.Error(err)
		return
	}
}
