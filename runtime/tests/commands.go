package tests

import (
	"github.com/go-primus/doma/runtime/core"
	"github.com/google/uuid"
)

const (
	CreateDemoCommand core.CommandType = "createdemo"
)

type CreateDemo struct {
	ID   uuid.UUID
	Name string
}

func (c CreateDemo) AggregateID() uuid.UUID            { return c.ID }
func (c CreateDemo) AggregateType() core.AggregateType { return TestAggregateRegisterType }
func (c CreateDemo) CommandType() core.CommandType     { return CreateDemoCommand }
