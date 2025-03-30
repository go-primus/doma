package tests

import (
	"context"
	"fmt"
	"time"

	"github.com/go-primus/doma/runtime/aggregatestore"
	"github.com/go-primus/doma/runtime/core"
	"github.com/google/uuid"
)

const (
	TestAggregateRegisterType      core.AggregateType = "demo"
	TestAggregateRegisterEmptyType core.AggregateType = ""
	TestAggregateRegisterTwiceType core.AggregateType = "TestAggregateRegisterTwice"
)

type TestAggregateRegister struct {
	*aggregatestore.AggregateBase
	id uuid.UUID
}

func newversion() aggregatestore.VersionedAggregate {
	return &TestAggregateRegister{}
}

var _ = core.Aggregate(&TestAggregateRegister{})

// NewInvitationAggregate creates a new InvitationAggregate with an ID.
func NewInvitationAggregate(id uuid.UUID) *TestAggregateRegister {
	return &TestAggregateRegister{
		AggregateBase: aggregatestore.NewAggregateBase(TestAggregateRegisterType, id),
	}
}

// func (a *TestAggregateRegister) EntityID() uuid.UUID { return a.id }

//	func (a *TestAggregateRegister) AggregateType() doma.AggregateType {
//		return TestAggregateRegisterType
//	}
func (a *TestAggregateRegister) HandleCommand(ctx context.Context, cmd core.Command) error {
	fmt.Println("handle command aggreate ", cmd)
	switch cmd := cmd.(type) {
	case *CreateDemo:
		a.AppendEvent(CreateDemoEvent, DemoCreateData{
			Name: cmd.Name,
		}, time.Now())
		return nil
	}
	return fmt.Errorf("invalid command %s, %s", cmd.AggregateType(), cmd.CommandType())
}

func (a *TestAggregateRegister) ApplyEvent(ctx context.Context, event core.Event) error {
	fmt.Println("apply event --", event)
	return nil
}
