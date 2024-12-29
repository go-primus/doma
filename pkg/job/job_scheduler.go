package job

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

// JobScheduler schedules jobs,
// all jobs within the same collection will run sequentially
const (
	collectionQueueCap = 64
	waitQueueCap       = 512
)

type jobQueue chan Job

type Scheduler struct {
	cancel context.CancelFunc
	wg     sync.WaitGroup

	processors *ConcurrentSet[int64] // Collections of having processor
	queues     map[int64]jobQueue    // CollectionID -> Queue
	waitQueue  jobQueue

	stopOne sync.Once
}

func NewScheduler() *Scheduler {
	return &Scheduler{
		processors: NewConcurrentSet[int64](),
		queues:     make(map[int64]jobQueue),
		waitQueue:  make(jobQueue, waitQueueCap),
	}
}

func (scheduler *Scheduler) Start() {
	ctx, cancel := context.WithCancel(context.Background())
	scheduler.cancel = cancel

	scheduler.wg.Add(1)
	go func() {
		defer scheduler.wg.Done()
		scheduler.schedule(ctx)
	}()
}

func (scheduler *Scheduler) Stop() {
	scheduler.stopOne.Do(func() {
		if scheduler.cancel != nil {
			scheduler.cancel()
		}

		scheduler.wg.Wait()
	})
}

func (scheduler *Scheduler) schedule(ctx context.Context) {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			slog.Info("JobManager stopped")
			for _, queue := range scheduler.queues {
				close(queue)
			}
			return
		case job := <-scheduler.waitQueue:
			queue, ok := scheduler.queues[job.CollectionID()]
			if !ok {
				queue = make(jobQueue, collectionQueueCap)
				scheduler.queues[job.CollectionID()] = queue
			}
			queue <- job
			scheduler.startProcessor(job.CollectionID(), queue)
		case <-ticker.C:
			for collection, queue := range scheduler.queues {
				if len(queue) > 0 {
					scheduler.startProcessor(collection, queue)
				} else {

					// Release resource if no job for the collection
					close(queue)
					delete(scheduler.queues, collection)
				}
			}
		}
	}
}

func (scheduler *Scheduler) Add(job Job) {
	scheduler.waitQueue <- job
}

func (scheduler *Scheduler) startProcessor(collection int64, queue jobQueue) {
	if !scheduler.processors.Insert(collection) {
		return
	}

	scheduler.wg.Add(1)
	go scheduler.processQueue(collection, queue)
}

// processQueue processes jobs in the given queue,
// it only processes jobs with the number of the length of queue at the time,
// to avoid leaking goroutines
func (scheduler *Scheduler) processQueue(collection int64, queue jobQueue) {
	defer scheduler.wg.Done()
	defer scheduler.processors.Remove(collection)

	len := len(queue)
	for i := 0; i < len; i++ {
		scheduler.process(<-queue)
	}
}

func (scheduler *Scheduler) process(job Job) {

	log := slog.With(slog.Int64("collectionID", job.CollectionID()))

	defer func() {
		log.Info("start to post-execute job")
		job.PostExecute()
		log.Info("job finished")
		job.Done()
	}()

	log.Info("start to pre-execute job")
	err := job.PreExecute()
	if err != nil {
		log.Warn("failed to pre-execute job", slog.Any("err", err))
		job.SetError(err)
		return
	}

	log.Info("start to execute job")
	err = job.Execute()
	if err != nil {
		log.Warn("failed to execute job", slog.Any("err", err))
		job.SetError(err)
	}
}
