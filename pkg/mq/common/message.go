package common

type Message interface {
	ID() string
	Topic() string
	Metadata() map[string]string
	Payload() []byte
}
