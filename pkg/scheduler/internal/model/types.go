package model

import "time"

type ITask interface {
	Wait() error
	Result() []byte
}

////////////////////

type TaskObject struct {
}

type TaskResource struct {
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

// ////////////////////
type TaskStatus struct {
	Status TaskState // pending、running、completed、failed
	Err    error     // 任务错误
	Result any       // 任务结果
	Synced bool      // 是否同步，最近同步时间
}

// /////////////////////////
type SyncMessage struct {
	Task     TaskMessage
	Status   TaskStatus
	Progress TaskProgress
}

// /////////
type Taskx struct {
	ID   string
	Type string
	// Version string

	Metadata TaskMetadata

	Spec     TaskSpec
	Status   TaskStatus
	Progress TaskProgress
}

type TaskMetadata struct {
	CreatedAt time.Time
	Priority  int32
	Labels    map[string]string
}

type TaskSpec struct {

	// 调度策略
	// 重试策略
	// 超时控制

}

type TaskProgress struct {
	Progress int // 百分比
	Success  int // 已完成
	Total    int
	//
	// CheckPoint string
	// Details string
}
