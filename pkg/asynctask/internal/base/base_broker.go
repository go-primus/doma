package base

import (
	"context"
	"time"
)

// Broker is a message broker that supports operations to manage task queues.
//
// See rdb.RDB as a reference implementation.
type Broker interface {
	// Ping() error
	// Close() error
	Enqueue(ctx context.Context, msg *TaskMessage) error
	EnqueueUnique(ctx context.Context, msg *TaskMessage, ttl time.Duration) error
	Dequeue(qnames ...string) (*TaskMessage, time.Time, error)
	// Done(ctx context.Context, msg *TaskMessage) error
	// MarkAsComplete(ctx context.Context, msg *TaskMessage) error
	// Requeue(ctx context.Context, msg *TaskMessage) error
	// Schedule(ctx context.Context, msg *TaskMessage, processAt time.Time) error
	// ScheduleUnique(ctx context.Context, msg *TaskMessage, processAt time.Time, ttl time.Duration) error
	// Retry(ctx context.Context, msg *TaskMessage, processAt time.Time, errMsg string, isFailure bool) error
	// Archive(ctx context.Context, msg *TaskMessage, errMsg string) error
	// ForwardIfReady(qnames ...string) error

	// // Group aggregation related methods
	// AddToGroup(ctx context.Context, msg *TaskMessage, gname string) error
	// AddToGroupUnique(ctx context.Context, msg *TaskMessage, groupKey string, ttl time.Duration) error
	// ListGroups(qname string) ([]string, error)
	// AggregationCheck(qname, gname string, t time.Time, gracePeriod, maxDelay time.Duration, maxSize int) (aggregationSetID string, err error)
	// ReadAggregationSet(qname, gname, aggregationSetID string) ([]*TaskMessage, time.Time, error)
	// DeleteAggregationSet(ctx context.Context, qname, gname, aggregationSetID string) error
	// ReclaimStaleAggregationSets(qname string) error

	// // Task retention related method
	// DeleteExpiredCompletedTasks(qname string) error

	// // Lease related methods
	// ListLeaseExpired(cutoff time.Time, qnames ...string) ([]*TaskMessage, error)
	// ExtendLease(qname string, ids ...string) (time.Time, error)

	// // State snapshot related methods
	// WriteServerState(info *ServerInfo, workers []*WorkerInfo, ttl time.Duration) error
	// ClearServerState(host string, pid int, serverID string) error

	// // Cancelation related methods
	// CancelationPubSub() (*redis.PubSub, error) // TODO: Need to decouple from redis to support other brokers
	// PublishCancelation(id string) error

	// WriteResult(qname, id string, data []byte) (n int, err error)
}
