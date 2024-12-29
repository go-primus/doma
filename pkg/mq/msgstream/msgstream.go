package msgstream

import "context"

type Timestamp uint64

// MsgPack represents a batch of msg in msgstream
type MsgPack struct {
	BeginTs        Timestamp
	EndTs          Timestamp
	Msgs           []TsMsg
	StartPositions []*MsgPosition
	EndPositions   []*MsgPosition
}

type MsgStream interface {
	AsProducer(ctx context.Context, channel string)
	Produce(context.Context, *MsgPack) error

	AsConsumer(ctx context.Context, channel string, subName string, position int) error
	Chan() <-chan *MsgPack
	// Seek(ctx)
	GetLatestMsgID() (string, error)
	CheckTopicValid() error

	EnableProduce(can bool)
}

type mqMsgStream struct{}

func NewMqMsgStream() (*mqMsgStream, error) {
	stream := &mqMsgStream{}

	return stream, nil
}
