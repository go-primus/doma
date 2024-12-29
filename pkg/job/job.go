package job

import "context"

type Job interface {
	MsgID() int64
	CollectionID() int64

	PreExecute() error
	Execute() error
	PostExecute() error

	Error() error
	SetError(err error)

	Done()
	Wait() error
}

type BaseJob struct {
	ctx          context.Context
	msgID        int64
	collectionID int64
	err          error
	doneCh       chan struct{}
}

func NewBaseJob(ctx context.Context, msgID, collectionID int64) *BaseJob {
	return &BaseJob{
		ctx:          ctx,
		msgID:        msgID,
		collectionID: collectionID,
		doneCh:       make(chan struct{}),
	}
}

func (job *BaseJob) MsgID() int64 {
	return job.msgID
}

func (job *BaseJob) CollectionID() int64 {
	return job.collectionID
}

func (job *BaseJob) Context() context.Context {
	return job.ctx
}

func (job *BaseJob) Error() error {
	return job.err
}

func (job *BaseJob) SetError(err error) {
	job.err = err
}

func (job *BaseJob) Done() {
	close(job.doneCh)
}

func (job *BaseJob) Wait() error {
	<-job.doneCh
	return job.err
}

func (job *BaseJob) PreExecute() error {
	return nil
}

func (job *BaseJob) PostExecute() {}
