package asynctask

import "errors"

var (
	// ErrNoProcessableTask indicates that there are no tasks ready to be processed.
	ErrNoProcessableTask = errors.New("no tasks are ready for processing")

	// ErrDuplicateTask indicates that another task with the same unique key holds the uniqueness lock.
	ErrDuplicateTask = errors.New("task already exists")

	// ErrTaskIdConflict indicates that another task with the same task ID already exist
	ErrTaskIdConflict = errors.New("task id conflicts with another task")
)
