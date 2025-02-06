package lifecycle

type Lifecycle[T any] interface {
	SetState(state T)
	GetState() T
}
