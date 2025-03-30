package core

type AggregateType string

func (at AggregateType) String() string {
	return string(at)
}

type Aggregate interface {
	Entity
	AggregateType() AggregateType
	CommandHandler
}
