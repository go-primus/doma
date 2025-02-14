package model

import "time"

// [Pending] -> [Running] <-> [Paused]
//     |           |
//     v           v
//  [Completed] [Failed]

type TaskState string

const (
	TaskState_Pending   TaskState = "pending"
	TaskState_Running   TaskState = "running"
	TaskState_Paused    TaskState = "paused"
	TaskState_Completed TaskState = "completed"
	TaskState_Cancelled TaskState = "cancelled"
	TaskState_Failed    TaskState = "failed"
)

////////
// task id, type, meta, spec, status, progress, result , checkpoint

type TaskMetadata struct {
	Queue     string // namespace
	CreatedAt time.Time
	Priority  int32

	Retry   int // 任务最大重试次数
	Timeout int64
	// Retention int64

	Labels map[string]string
}

type TaskSpec struct {

	// 调度策略
	// 重试策略
	// 超时控制

}

type TaskStatus struct {
	Status     TaskState // pending、running、completed、failed
	Err        error     // 任务错误
	Result     any       // 任务结果
	CheckPoint any       // checkpoint
	Synced     bool      // 是否同步，最近同步时间
}

type TaskProgress struct {
	Progress int // 百分比
	Success  int // 已完成
	Total    int
	//
	// CheckPoint string
	// Details string
}

////////////////////

type TaskObject struct {
}

type TaskResource struct {
}
