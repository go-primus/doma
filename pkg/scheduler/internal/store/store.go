package store

import (
	"fmt"
	"sync"

	"github.com/go-primus/doma/pkg/scheduler/internal/model"
)

type TaskStore interface {
	AddTask(task model.TaskMessage)
	PauseTask(taskId string)
	ResumeTask(taskId string)
	UpdateStatus(taskId string, status model.TaskStatus)
	UpdateProgress(taskId string, progress model.TaskProgress)

	SaveCheckPoint(taskId string, checkPoint any)
	LoadCheckPoint(taskId string) any

	GetStatus(taskId string) model.TaskStatus
	GetTask(taskId string) *TaskItem

	MarkSynced(taskId string)
	ListSyncTasks() []*TaskItem
}

var _ TaskStore = (*taskStore)(nil)

type TaskItem struct {
	Msg        model.TaskMessage
	Status     model.TaskStatus
	Progress   model.TaskProgress
	Checkpoint any
	Synced     bool
}

type taskStore struct {
	tasks map[string]*TaskItem
	mu    sync.Mutex
}

func NewTaskStore() *taskStore {
	return &taskStore{
		tasks: make(map[string]*TaskItem),
	}
}

// LoadCheckPoint implements TaskStore.
func (t *taskStore) LoadCheckPoint(taskId string) any {
	t.mu.Lock()
	defer t.mu.Unlock()

	task, ok := t.tasks[taskId]
	if !ok {
		return nil
	}

	return task.Checkpoint
}

// SaveCheckPoint implements TaskStore.
func (t *taskStore) SaveCheckPoint(taskId string, checkPoint any) {
	t.mu.Lock()
	defer t.mu.Unlock()

	task, ok := t.tasks[taskId]
	if !ok {
		return
	}

	task.Checkpoint = checkPoint

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
func (t *taskStore) UpdateProgress(taskId string, progress model.TaskProgress) {
	t.mu.Lock()
	defer t.mu.Unlock()

	item, ok := t.tasks[taskId]
	if !ok {
		return
	}

	item.Progress = progress
	item.Synced = false
}

// MarkSynced implements TaskStore.
func (t *taskStore) MarkSynced(taskId string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	task, ok := t.tasks[taskId]
	if !ok {
		return
	}

	task.Synced = true
}

// GetUnSyncTask implements TaskStore.
func (t *taskStore) ListSyncTasks() []*TaskItem {
	t.mu.Lock()
	defer t.mu.Unlock()
	tasks := []*TaskItem{}
	for _, task := range t.tasks {
		if !task.Synced {
			tasks = append(tasks, task)
		}
	}
	return tasks
}

func (t *taskStore) AddTask(task model.TaskMessage) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.tasks[task.ID] = &TaskItem{
		Msg: task,
		Status: model.TaskStatus{
			Status: model.TaskState_Running,
		},
		Synced: false,
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

	item.Status.Status = "paused"
	item.Synced = false
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

	item.Status.Status = "running"
	item.Synced = false
	t.tasks[taskId] = item
}

// GetStatus implements TaskStore.
func (t *taskStore) GetStatus(taskId string) model.TaskStatus {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.tasks[taskId].Status
}

// UpdateStatus implements TaskStore.
func (t *taskStore) UpdateStatus(taskId string, status model.TaskStatus) {
	t.mu.Lock()
	defer t.mu.Unlock()
	fmt.Println("update task :", taskId, "---", status.Status)
	item, ok := t.tasks[taskId]
	if !ok {
		return
	}
	item.Status = status
	item.Synced = false
	t.tasks[taskId] = item
}
