package base

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/go-primus/doma/pkg/asynctask/internal/errors"
)

// Version of asynq library and CLI.
const Version = "0.24.1"

// DefaultQueueName is the queue name used if none are specified by user.
const DefaultQueueName = "default"

// DefaultQueue is the redis key for the default queue.
var DefaultQueue = PendingKey(DefaultQueueName)

// Global Redis keys.
const (
	AllServers    = "asynq:servers"    // ZSET
	AllWorkers    = "asynq:workers"    // ZSET
	AllSchedulers = "asynq:schedulers" // ZSET
	AllQueues     = "asynq:queues"     // SET
	CancelChannel = "asynq:cancel"     // PubSub channel
)

// TaskState denotes the state of a task.
type TaskState int

const (
	TaskStateActive TaskState = iota + 1
	TaskStatePending
	TaskStateScheduled
	TaskStateRetry
	TaskStateArchived
	TaskStateCompleted
	TaskStateAggregating // describes a state where task is waiting in a group to be aggregated
)

func (s TaskState) String() string {
	switch s {
	case TaskStateActive:
		return "active"
	case TaskStatePending:
		return "pending"
	case TaskStateScheduled:
		return "scheduled"
	case TaskStateRetry:
		return "retry"
	case TaskStateArchived:
		return "archived"
	case TaskStateCompleted:
		return "completed"
	case TaskStateAggregating:
		return "aggregating"
	}
	panic(fmt.Sprintf("internal error: unknown task state %d", s))
}

func TaskStateFromString(s string) (TaskState, error) {
	switch s {
	case "active":
		return TaskStateActive, nil
	case "pending":
		return TaskStatePending, nil
	case "scheduled":
		return TaskStateScheduled, nil
	case "retry":
		return TaskStateRetry, nil
	case "archived":
		return TaskStateArchived, nil
	case "completed":
		return TaskStateCompleted, nil
	case "aggregating":
		return TaskStateAggregating, nil
	}
	return 0, errors.E(errors.FailedPrecondition, fmt.Sprintf("%q is not supported task state", s))
}

// ValidateQueueName validates a given qname to be used as a queue name.
// Returns nil if valid, otherwise returns non-nil error.
func ValidateQueueName(qname string) error {
	if len(strings.TrimSpace(qname)) == 0 {
		return fmt.Errorf("queue name must contain one or more characters")
	}
	return nil
}

// QueueKeyPrefix returns a prefix for all keys in the given queue.
func QueueKeyPrefix(qname string) string {
	return fmt.Sprintf("asynq:{%s}:", qname)
}

// TaskKeyPrefix returns a prefix for task key.
func TaskKeyPrefix(qname string) string {
	return fmt.Sprintf("%st:", QueueKeyPrefix(qname))
}

// TaskKey returns a redis key for the given task message.
func TaskKey(qname, id string) string {
	return fmt.Sprintf("%s%s", TaskKeyPrefix(qname), id)
}

// PendingKey returns a redis key for the given queue name.
func PendingKey(qname string) string {
	return fmt.Sprintf("%spending", QueueKeyPrefix(qname))
}

// ActiveKey returns a redis key for the active tasks.
func ActiveKey(qname string) string {
	return fmt.Sprintf("%sactive", QueueKeyPrefix(qname))
}

// ScheduledKey returns a redis key for the scheduled tasks.
func ScheduledKey(qname string) string {
	return fmt.Sprintf("%sscheduled", QueueKeyPrefix(qname))
}

// RetryKey returns a redis key for the retry tasks.
func RetryKey(qname string) string {
	return fmt.Sprintf("%sretry", QueueKeyPrefix(qname))
}

// ArchivedKey returns a redis key for the archived tasks.
func ArchivedKey(qname string) string {
	return fmt.Sprintf("%sarchived", QueueKeyPrefix(qname))
}

// LeaseKey returns a redis key for the lease.
func LeaseKey(qname string) string {
	return fmt.Sprintf("%slease", QueueKeyPrefix(qname))
}

func CompletedKey(qname string) string {
	return fmt.Sprintf("%scompleted", QueueKeyPrefix(qname))
}

// PausedKey returns a redis key to indicate that the given queue is paused.
func PausedKey(qname string) string {
	return fmt.Sprintf("%spaused", QueueKeyPrefix(qname))
}

// ProcessedTotalKey returns a redis key for total processed count for the given queue.
func ProcessedTotalKey(qname string) string {
	return fmt.Sprintf("%sprocessed", QueueKeyPrefix(qname))
}

// FailedTotalKey returns a redis key for total failure count for the given queue.
func FailedTotalKey(qname string) string {
	return fmt.Sprintf("%sfailed", QueueKeyPrefix(qname))
}

// ProcessedKey returns a redis key for processed count for the given day for the queue.
func ProcessedKey(qname string, t time.Time) string {
	return fmt.Sprintf("%sprocessed:%s", QueueKeyPrefix(qname), t.UTC().Format("2006-01-02"))
}

// FailedKey returns a redis key for failure count for the given day for the queue.
func FailedKey(qname string, t time.Time) string {
	return fmt.Sprintf("%sfailed:%s", QueueKeyPrefix(qname), t.UTC().Format("2006-01-02"))
}

// ServerInfoKey returns a redis key for process info.
func ServerInfoKey(hostname string, pid int, serverID string) string {
	return fmt.Sprintf("asynq:servers:{%s:%d:%s}", hostname, pid, serverID)
}

// WorkersKey returns a redis key for the workers given hostname, pid, and server ID.
func WorkersKey(hostname string, pid int, serverID string) string {
	return fmt.Sprintf("asynq:workers:{%s:%d:%s}", hostname, pid, serverID)
}

// SchedulerEntriesKey returns a redis key for the scheduler entries given scheduler ID.
func SchedulerEntriesKey(schedulerID string) string {
	return fmt.Sprintf("asynq:schedulers:{%s}", schedulerID)
}

// SchedulerHistoryKey returns a redis key for the scheduler's history for the given entry.
func SchedulerHistoryKey(entryID string) string {
	return fmt.Sprintf("asynq:scheduler_history:%s", entryID)
}

// UniqueKey returns a redis key with the given type, payload, and queue name.
func UniqueKey(qname, tasktype string, payload []byte) string {
	if payload == nil {
		return fmt.Sprintf("%sunique:%s:", QueueKeyPrefix(qname), tasktype)
	}
	checksum := md5.Sum(payload)
	return fmt.Sprintf("%sunique:%s:%s", QueueKeyPrefix(qname), tasktype, hex.EncodeToString(checksum[:]))
}

// GroupKeyPrefix returns a prefix for group key.
func GroupKeyPrefix(qname string) string {
	return fmt.Sprintf("%sg:", QueueKeyPrefix(qname))
}

// GroupKey returns a redis key used to group tasks belong in the same group.
func GroupKey(qname, gkey string) string {
	return fmt.Sprintf("%s%s", GroupKeyPrefix(qname), gkey)
}

// AggregationSetKey returns a redis key used for an aggregation set.
func AggregationSetKey(qname, gname, setID string) string {
	return fmt.Sprintf("%s:%s", GroupKey(qname, gname), setID)
}

// AllGroups return a redis key used to store all group keys used in a given queue.
func AllGroups(qname string) string {
	return fmt.Sprintf("%sgroups", QueueKeyPrefix(qname))
}

// AllAggregationSets returns a redis key used to store all aggregation sets (set of tasks staged to be aggregated)
// in a given queue.
func AllAggregationSets(qname string) string {
	return fmt.Sprintf("%saggregation_sets", QueueKeyPrefix(qname))
}

// TaskMessage is the internal representation of a task with additional metadata fields.
// Serialized data of this type gets written to redis.
type TaskMessage struct {
	// Type indicates the kind of the task to be performed.
	Type string

	// Payload holds data needed to process the task.
	Payload []byte

	// ID is a unique identifier for each task.
	ID string

	// Queue is a name this message should be enqueued to.
	Queue string

	// Retry is the max number of retry for this task.
	Retry int

	// Retried is the number of times we've retried this task so far.
	Retried int

	// ErrorMsg holds the error message from the last failure.
	ErrorMsg string

	// Time of last failure in Unix time,
	// the number of seconds elapsed since January 1, 1970 UTC.
	//
	// Use zero to indicate no last failure
	LastFailedAt int64

	// Timeout specifies timeout in seconds.
	// If task processing doesn't complete within the timeout, the task will be retried
	// if retry count is remaining. Otherwise it will be moved to the archive.
	//
	// Use zero to indicate no timeout.
	Timeout int64

	// Deadline specifies the deadline for the task in Unix time,
	// the number of seconds elapsed since January 1, 1970 UTC.
	// If task processing doesn't complete before the deadline, the task will be retried
	// if retry count is remaining. Otherwise it will be moved to the archive.
	//
	// Use zero to indicate no deadline.
	Deadline int64

	// UniqueKey holds the redis key used for uniqueness lock for this task.
	//
	// Empty string indicates that no uniqueness lock was used.
	UniqueKey string

	// GroupKey holds the group key used for task aggregation.
	//
	// Empty string indicates no aggregation is used for this task.
	GroupKey string

	// Retention specifies the number of seconds the task should be retained after completion.
	Retention int64

	// CompletedAt is the time the task was processed successfully in Unix time,
	// the number of seconds elapsed since January 1, 1970 UTC.
	//
	// Use zero to indicate no value.
	CompletedAt int64
}
