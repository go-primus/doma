package scheduler

type Client interface {
	SubmitTask() (TaskInfo, error)
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
