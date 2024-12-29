package channels

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-primus/doma/pkg/mq/msgstream"
)

type ChannelsMgr interface {
	GetChannel() string
	GetOrCreateStream(ctx context.Context) (msgstream.MsgStream, error)
	RemoveStream()
}

var _ ChannelsMgr = &channelsMgrImpl{}

type streamInfo struct {
	channelInfo string
	stream      msgstream.MsgStream
}

type channelsMgrImpl struct {
	mu   sync.RWMutex
	info streamInfo

	msgStreamFactory msgstream.Factory
}

// GetChannel implements ChannelsMgr.
func (c *channelsMgrImpl) GetChannel() string {
	return ""
}

// GetOrCreateStream implements ChannelsMgr.
func (mgr *channelsMgrImpl) GetOrCreateStream(ctx context.Context) (msgstream.MsgStream, error) {

	if stream, err := mgr.localGetStream(); err == nil {
		return stream, nil
	}

	return mgr.createMsgStream(ctx)
}

// RemoveStream implements ChannelsMgr.
func (c *channelsMgrImpl) RemoveStream() {
	panic("unimplemented")
}

// //internal
func (mgr *channelsMgrImpl) localGetStream() (msgstream.MsgStream, error) {
	mgr.mu.RLock()
	defer mgr.mu.Unlock()

	if mgr.info.stream != nil {
		return mgr.info.stream, nil
	}

	return nil, fmt.Errorf("stream not found")
}

func (mgr *channelsMgrImpl) createMsgStream(ctx context.Context) (msgstream.MsgStream, error) {
	mgr.mu.RLock()
	if mgr.info.stream != nil {
		mgr.mu.Unlock()
		return mgr.info.stream, nil
	}

	mgr.mu.RUnlock()

	stream, err := createStream(ctx, mgr.msgStreamFactory)
	if err != nil {
		return nil, err
	}

	mgr.mu.Lock()
	defer mgr.mu.Unlock()
	mgr.info.channelInfo = ""
	mgr.info.stream = stream
	return mgr.info.stream, nil

}

func createStream(ctx context.Context, factory msgstream.Factory) (msgstream.MsgStream, error) {
	var stream msgstream.MsgStream
	var err error

	stream, err = factory.NewMsgStream(context.Background())
	if err != nil {
		return nil, err
	}

	stream.AsProducer(ctx, pchans)
	if repack != nil {
		stream.SetRepackFunc(repack)
	}
	return stream, nil
}
