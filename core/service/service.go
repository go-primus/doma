package service

import "context"

type Service interface {
	// start service
	Start(ctx context.Context) error
	// stop service
	Stop(ctx context.Context) error
}
