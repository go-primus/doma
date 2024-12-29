package msgstream

import "context"

type Factory interface {
	NewMsgStream(ctx context.Context) (MsgStream, error)
}

func NewMqFactory() Factory {

	return &CommonFactory{}
}

var _ Factory = &CommonFactory{}

type CommonFactory struct{}

// NewMsgStream implements Factory.
func (c *CommonFactory) NewMsgStream(ctx context.Context) (MsgStream, error) {

	return NewMqMsgStream()
}
