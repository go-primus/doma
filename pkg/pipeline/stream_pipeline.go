package pipeline

import "context"

type StreamPipeline interface {
	Pipeline
	ConsumeStream(ctx context.Context) error
	Status() string
}
