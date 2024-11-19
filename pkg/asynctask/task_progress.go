package asynctask

type TaskProgress interface {
	UpdateProgress(success, total int32)
}

type taskProgress struct {
	w *ResultWriter
}

// UpdateProgress implements TaskProgress.
func (t *taskProgress) UpdateProgress(success int32, total int32) {

	t.w.Write([]byte(""))
}

func NewTaskProgress() TaskProgress {
	return &taskProgress{}
}
