package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/primus/primus/doma"
	"github.com/primus/primus/doma/aggregatestore"
	memoryEventStore "github.com/primus/primus/doma/eventstore/memory"
)

type DomaDemoService interface {
	CreateDemo(uuid.UUID, string) error
}

type demoDemoServiceImpl struct {
	t     doma.AggregateType
	store doma.AggregateStore
}

func (d *demoDemoServiceImpl) handleCommand(ctx context.Context, cmd doma.Command) error {
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

func NewDemoService(estore doma.EventStore) DomaDemoService {
	var store doma.AggregateStore
	store, _ = aggregatestore.NewAggregateStore(estore)
	return &demoDemoServiceImpl{
		t:     TestAggregateRegisterType,
		store: store,
	}
}

func TestDoma(t *testing.T) {

	// 注册聚合

	doma.RegisterAggregate(func(u uuid.UUID) doma.Aggregate {
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

	err = s.CreateDemo(id, "name")
	if err != nil {
		t.Error(err)
		return
	}
}
