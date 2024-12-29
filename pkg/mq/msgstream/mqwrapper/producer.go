package mqwrapper

import (
	"context"

	"github.com/go-primus/doma/pkg/mq/common"
)

type Producer interface {
	Send(ctx context.Context, message common.Message) (string, error)
	Close()
}
