package scheduler

import (
	"fmt"
	"sync"
)

type TaskStore interface {
	AddTask(task TaskMessage)
	PauseTask(taskId string)
	ResumeTask(taskId string)
	UpdateStatus(taskId string, status TaskStatus)
	GetStatus(taskId string) TaskStatus

	MarkSynced(taskId string)
	ListSyncStatus() []TaskStatus
}

var _ TaskStore = (*taskStore)(nil)

type TaskItem struct {
	msg    TaskMessage
	status TaskStatus
	synced bool
}

type taskStore struct {
	tasks map[string]*TaskItem
	mu    sync.Mutex
}

// MarkSynced implements TaskStore.
func (t *taskStore) MarkSynced(taskId string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	task, ok := t.tasks[taskId]
	if !ok {
		return
	}

	task.synced = true
}

// GetUnSyncTask implements TaskStore.
func (t *taskStore) ListSyncStatus() []TaskStatus {
	t.mu.Lock()
	defer t.mu.Unlock()
	tasks := []TaskStatus{}
	for _, task := range t.tasks {
		if !task.synced {
			tasks = append(tasks, task.status)
		}
	}
	return tasks
}

func (t *taskStore) AddTask(task TaskMessage) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.tasks[task.ID] = &TaskItem{
		msg: task,
		status: TaskStatus{
			ID:     task.ID,
			Status: "running",
		},
		synced: false,
	}
}

// PauseTask implements TaskStore.
func (t *taskStore) PauseTask(taskId string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	item, ok := t.tasks[taskId]
	if !ok {
		return
	}

	item.status.Status = "paused"
	item.synced = false
	t.tasks[taskId] = item

}

// ResumeTask implements TaskStore.
func (t *taskStore) ResumeTask(taskId string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	item, ok := t.tasks[taskId]
	if !ok {
		return
	}

	item.status.Status = "running"
	item.synced = false
	t.tasks[taskId] = item
}

// GetStatus implements TaskStore.
func (t *taskStore) GetStatus(taskId string) TaskStatus {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.tasks[taskId].status
}

// UpdateStatus implements TaskStore.
func (t *taskStore) UpdateStatus(taskId string, status TaskStatus) {
	t.mu.Lock()
	defer t.mu.Unlock()
	fmt.Println("update task :", status.ID, "---", status.Status)
	item, ok := t.tasks[taskId]
	if !ok {
		return
	}
	item.status = status
	item.synced = false
	t.tasks[taskId] = item
}
