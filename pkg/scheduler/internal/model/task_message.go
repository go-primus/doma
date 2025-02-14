package model

type TaskMessage struct {
	ID   string // unique identifier
	Type string // task kind

	Metadata TaskMetadata

	Payload any
}
