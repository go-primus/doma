package mqwrapper

import "context"

type Client interface {
	CreateProducer(ctx context.Context) (Producer, error)
	Subscribe(ctx context.Context) (Consumer, error)

	FirstMsgID() (string, error)
	Close()
}
