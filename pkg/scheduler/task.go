package scheduler

type Task interface {
	Wait() error
	Result() []byte
}
