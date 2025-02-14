package store

import "github.com/go-primus/doma/pkg/scheduler/internal/model"

type TaskStore interface {
	AddTask(task model.TaskMessage)
	PauseTask(taskId string)
	ResumeTask(taskId string)

	UpdateStatus(taskId string, status model.TaskStatus)
	UpdateProgress(taskId string, progress model.TaskProgress)

	// check point
	SaveCheckPoint(taskId string, checkPoint any)
	LoadCheckPoint(taskId string) any

	GetStatus(taskId string) model.TaskStatus
	GetTask(taskId string) *TaskItem

	// syncing
	MarkSynced(taskId string)
	ListSyncTasks() []*TaskItem
}

type TaskItem struct {
	Msg        model.TaskMessage
	Status     model.TaskStatus
	Progress   model.TaskProgress
	Checkpoint any
	Synced     bool
}
