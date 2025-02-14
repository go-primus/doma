package model

// /////////
type Task struct {
	ID       string
	Type     string
	Metadata TaskMetadata
	// Version string
	Payload any
}
