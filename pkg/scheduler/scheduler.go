package scheduler

type Scheduler interface {
	Start() error
	Stop() error

	Add(task Task) error

	SubmitTask(task Task) (string, error)

	// 暂停、取消、重试任务
	PauseTask(id string) error
	ResumeTask(id string) error
	DeleteTask(id string) error
	RetryTask(id string) error

	//
	GetTask(id string) TaskStatus
}

type SchedulePolicy interface {
	Push(task Task) (int, error)
	Pop() Task
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
