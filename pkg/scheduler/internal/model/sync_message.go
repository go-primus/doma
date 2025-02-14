package model

// /////////////////////////
type SyncMessage struct {
	Task     TaskMessage
	Status   TaskStatus
	Progress TaskProgress
}
