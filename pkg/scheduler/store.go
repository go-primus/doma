package scheduler

type TaskStore interface {
	UpdateStatus(status TaskStatus)
	GetStatus(taskId string) TaskStatus
}

var _ TaskStore = (*taskStore)(nil)

type taskStore struct {
	tasks map[string]TaskStatus
}

// GetStatus implements TaskStore.
func (t *taskStore) GetStatus(taskId string) TaskStatus {
	return t.tasks[taskId]
}

// UpdateStatus implements TaskStore.
func (t *taskStore) UpdateStatus(status TaskStatus) {
	t.tasks[status.ID] = status
}
