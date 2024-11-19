package asynctask

import "context"

// type

// type JobProgressNotify func(api.JobProgressInfo)
type TaskCallback interface {
	OnProgress()
	OnError()
	OnComplete()
}

type TaskHandle func(ctx context.Context) error

type AsyncTask struct {
	name   string
	handle TaskHandle
}

func NewAsyncTask(ctx context.Context, name string, handle TaskHandle) AsyncTask {
	return AsyncTask{
		name:   name,
		handle: handle,
	}
}

func (task *AsyncTask) Execute() {

	ctx := context.Background()
	// dispatch
	task.handle(ctx)
}

func (task *AsyncTask) dispatch() {

}

// ///

type TaskStore interface {
	Dispatch()
	Commit()
}

// ///
type ProgressTask struct {
	AsyncTask
	//

	Progress int `json:"progress"` // job整体进度，百分制

	Success int `json:"success"` // 已成功完成数
	Total   int `json:"total"`   // 总共要完成的数

	Desc any `json:"desc"` // job自定义信息

}

func (task *ProgressTask) publishProgress(progress int) {

}
