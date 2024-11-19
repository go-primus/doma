package eventbus

type nullBus struct {
}

func NewNullBus() EventBus {
	return &nullBus{}
}

func (b *nullBus) Publish(topic string, event any) error {

	return nil
}
func (b *nullBus) Subscribe(topic string, fn EventHandler) {

}
func (b *nullBus) UnSubscribe(topic string, fn EventHandler) {

}
