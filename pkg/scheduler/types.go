package scheduler

type TaskMessage struct {
	ID      string
	Type    string
	Payload any
}

type Task struct {
	ID      string
	Type    string
	Payload any
}

type TaskCommand struct {
	ID      string
	Command string
}

// 成功
// 进行中

type TaskStatus struct {
	ID       string
	Status   string // pending、running、completed、failed
	Progress int
	Err      string
}

// /////////////////////////
type SyncMessage struct {
	task   TaskMessage
	status TaskStatus
}
