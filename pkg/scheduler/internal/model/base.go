package model

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
	TaskState_Failed    TaskState = "failed"
)

//////
