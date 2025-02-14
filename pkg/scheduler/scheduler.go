package scheduler

import "github.com/go-primus/doma/pkg/scheduler/internal/model"

type Scheduler interface {
	Start() error
	Stop() error

	Add(task model.Task) error

	SubmitTask(task model.Task) (string, error)

	// 暂停、取消、重试任务
	PauseTask(id string) error
	ResumeTask(id string) error
	DeleteTask(id string) error
	RetryTask(id string) error

	//
	GetTask(id string) model.TaskStatus
}

type SchedulePolicy interface {
	Push(task model.Task) (int, error)
	Pop() model.Task
	Len() int
}

func NewScheduler(policyName string) Scheduler {
	switch policyName {
	case "":
		fallthrough
	default:
	}
	return nil
}
