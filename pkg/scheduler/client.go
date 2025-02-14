package scheduler

import "github.com/go-primus/doma/pkg/scheduler/internal/model"

type Client interface {
	SubmitTask() (model.TaskInfo, error)
	//
	DeleteTask()
	//
	PauseTask()
	ResumeTask()
	//
	GetTaskStatus()
	GetTaskInfo()
	ListTasks()
}

type TaskResult interface {
	WaitToFinish()
	Notify(err error)
}
