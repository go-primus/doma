package msgstream

import (
	"context"
	"testing"
)

func TestMq(t *testing.T) {

	f := NewMqFactory()
	stream, _ := f.NewMsgStream(context.Background())
}
