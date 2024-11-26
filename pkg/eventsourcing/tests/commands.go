package tests

import (
	"github.com/google/uuid"
	"github.com/primus/primus/doma"
)

const (
	CreateDemoCommand doma.CommandType = "createdemo"
)

type CreateDemo struct {
	ID   uuid.UUID
	Name string
}

func (c CreateDemo) AggregateID() uuid.UUID            { return c.ID }
func (c CreateDemo) AggregateType() doma.AggregateType { return TestAggregateRegisterType }
func (c CreateDemo) CommandType() doma.CommandType     { return CreateDemoCommand }
