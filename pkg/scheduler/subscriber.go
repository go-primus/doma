package scheduler

import (
	"github.com/go-primus/doma/pkg/scheduler/internal/cancelation"
	"github.com/go-primus/doma/pkg/scheduler/internal/model"
)

type subscriber struct {
	done chan struct{}

	cancel <-chan model.TaskCommand

	cancelations *cancelation.Cancelations
}

type subscriberParams struct {
	cancel       <-chan model.TaskCommand
	cancelations *cancelation.Cancelations
}

func newSubscriber(params subscriberParams) *subscriber {
	return &subscriber{
		done:         make(chan struct{}),
		cancelations: params.cancelations,
		cancel:       params.cancel,
	}
}

func (s *subscriber) Shutdown() {
	s.done <- struct{}{}
}

func (s *subscriber) Start() {
	go func() {
		for {
			select {
			case <-s.done:
				return
			case msg := <-s.cancel:
				cancel, ok := s.cancelations.Get(msg.ID)
				if ok {
					cancel()
				}
			}
		}
	}()
}
