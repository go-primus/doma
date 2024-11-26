package tests

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/primus/primus/doma"
	"github.com/primus/primus/doma/aggregatestore"
)

const (
	TestAggregateRegisterType      doma.AggregateType = "demo"
	TestAggregateRegisterEmptyType doma.AggregateType = ""
	TestAggregateRegisterTwiceType doma.AggregateType = "TestAggregateRegisterTwice"
)

type TestAggregateRegister struct {
	*aggregatestore.AggregateBase
	id uuid.UUID
}

func newversion() aggregatestore.VersionedAggregate {
	return &TestAggregateRegister{}
}

var _ = doma.Aggregate(&TestAggregateRegister{})

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
func (a *TestAggregateRegister) HandleCommand(ctx context.Context, cmd doma.Command) error {
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

func (a *TestAggregateRegister) ApplyEvent(ctx context.Context, event doma.Event) error {
	fmt.Println("apply event ", event)
	return nil
}
