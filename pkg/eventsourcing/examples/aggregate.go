package examples

import "github.com/primus/primus/doma/aggregatestore"

type UserRepo interface {
}

type UserAggregate struct {
	*aggregatestore.AggregateBase

	//
	userRepo UserRepo
}
