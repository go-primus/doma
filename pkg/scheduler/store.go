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
	UpdateProgress(taskId string, progress TaskProgress)

	SaveCheckPoint(taskId string, checkPoint any)
	LoadCheckPoint(taskId string) any

	GetStatus(taskId string) TaskStatus
	GetTask(taskId string) *TaskItem

	MarkSynced(taskId string)
	ListSyncTasks() []*TaskItem
}

var _ TaskStore = (*taskStore)(nil)

type TaskItem struct {
	msg        TaskMessage
	status     TaskStatus
	progress   TaskProgress
	checkpoint any
	synced     bool
}

type taskStore struct {
	tasks map[string]*TaskItem
	mu    sync.Mutex
}

// LoadCheckPoint implements TaskStore.
func (t *taskStore) LoadCheckPoint(taskId string) any {
	t.mu.Lock()
	defer t.mu.Unlock()

	task, ok := t.tasks[taskId]
	if !ok {
		return nil
	}

	return task.checkpoint
}

// SaveCheckPoint implements TaskStore.
func (t *taskStore) SaveCheckPoint(taskId string, checkPoint any) {
	t.mu.Lock()
	defer t.mu.Unlock()

	task, ok := t.tasks[taskId]
	if !ok {
		return
	}

	task.checkpoint = checkPoint

}

// GetTask implements TaskStore.
func (t *taskStore) GetTask(taskId string) *TaskItem {
	t.mu.Lock()
	defer t.mu.Unlock()

	task, ok := t.tasks[taskId]
	if !ok {
		return nil
	}

	return task
}

// UpdateProgress implements TaskStore.
func (t *taskStore) UpdateProgress(taskId string, progress TaskProgress) {
	t.mu.Lock()
	defer t.mu.Unlock()

	item, ok := t.tasks[taskId]
	if !ok {
		return
	}

	item.progress = progress
	item.synced = false
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
func (t *taskStore) ListSyncTasks() []*TaskItem {
	t.mu.Lock()
	defer t.mu.Unlock()
	tasks := []*TaskItem{}
	for _, task := range t.tasks {
		if !task.synced {
			tasks = append(tasks, task)
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
			Status: TaskState_Running,
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
	fmt.Println("update task :", taskId, "---", status.Status)
	item, ok := t.tasks[taskId]
	if !ok {
		return
	}
	item.status = status
	item.synced = false
	t.tasks[taskId] = item
}
