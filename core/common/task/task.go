package task

import "github.com/google/uuid"

type Task struct {
	TaskId   string `json:"task_id,omitempty"`
	TaskType string `json:"task_type,omitempty"`
	Payload  any    `json:"payload,omitempty"`

	//
	Metadata *TaskMetadata `json:"metadata,omitempty"`
}

func NewTask(taskType string, payload any) Task {

	return Task{
		TaskId:   uuid.NewString(),
		TaskType: taskType,
		Payload:  payload,
	}
}
