package mqwrapper

import "github.com/go-primus/doma/pkg/mq/common"

type Consumer interface {
	Sbuscription() string
	Chan() <-chan common.Message
	Seek(string, bool) error
	Ack(string)
	Close()
	GetLatestMsgID() (string, error)
	CheckTopicValid(channel string) error
}
