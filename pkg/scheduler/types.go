package scheduler

type Task struct {
	ID      string
	Type    string
	Payload any
}

// 成功
// 进行中

type TaskStatus struct {
	ID       string
	Status   string // pending、running、completed、failed
	Progress int
	Err      string
}

// ////////////////////////
type TaskCtx interface {
	UpdateProgress(progress int32)
}
