package asynctask

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/primus/primus/pkg/logger"
)

type Handler interface {
	ProcessTask(context.Context, *Task) error
}

type HandlerFunc func(context.Context, *Task) error

func (fn HandlerFunc) ProcessTask(ctx context.Context, task *Task) error {
	return fn(ctx, task)
}

//

// ////////
type TaskMux struct {
	// s     Service
	_name string
	mu    sync.RWMutex
	m     map[string]muxEntry
	es    []muxEntry
	mws   []MiddlewareFunc
}

type muxEntry struct {
	h       Handler
	pattern string
}

type MiddlewareFunc func(Handler) Handler

// //

func NewTaskMux(muxName string) *TaskMux {
	mux := new(TaskMux)
	mux._name = muxName
	return mux
}

// func NewServiceEventMux(s Service) *EventMux {
// 	mux := NewEventMux(s.Name())
// 	mux.s = s
// 	return mux
// }

func (mux *TaskMux) name() string {

	if mux._name != "" {
		return mux._name
	}

	// if mux.s == nil {
	// 	return ""
	// }
	// return mux.s.Name()
	return "unknown"
}

func (mux *TaskMux) Use(mws ...MiddlewareFunc) {
	mux.mu.Lock()
	defer mux.mu.Unlock()

	mux.mws = append(mux.mws, mws...)
}

func (mux *TaskMux) ProcessTask(ctx context.Context, task *Task) error {

	h, pattern := mux.Handler(task)
	if pattern == "" {
		return h.ProcessTask(ctx, task)
	}
	logger.L().Debug("=======[", mux.name(), "] ---- handle task: ", task.Type(), " ==============")
	err := h.ProcessTask(ctx, task)
	if err != nil {
		logger.L().Error("=======[", mux.name(), "] ---- handle task end:", err)
		return err
	}
	logger.L().Debug("=======[", mux.name(), "] ---- handle task end: ok ==============")

	return nil
}

func (mux *TaskMux) Handler(t *Task) (h Handler, pattern string) {
	mux.mu.RLock()
	defer mux.mu.RUnlock()

	h, pattern = mux.match(t.Type())
	if h == nil {
		h, pattern = NotFoundHandler(), ""
	}
	for i := len(mux.mws) - 1; i >= 0; i-- {
		h = mux.mws[i](h)
	}
	return h, pattern
}

// Find a handler on a handler map given a typename string.
// Most-specific (longest) pattern wins.
func (mux *TaskMux) match(typename string) (h Handler, pattern string) {
	// Check for exact match first.
	v, ok := mux.m[typename]
	if ok {
		return v.h, v.pattern
	}

	// Check for longest valid match.
	// mux.es contains all patterns from longest to shortest.
	for _, e := range mux.es {
		if strings.HasPrefix(typename, e.pattern) {
			return e.h, e.pattern
		}
	}
	return nil, ""

}

func (mux *TaskMux) Handle(pattern string, handler Handler) {
	mux.mu.Lock()
	defer mux.mu.Unlock()

	if strings.TrimSpace(pattern) == "" {
		return
	}

	if handler == nil {
		return
	}

	if _, exist := mux.m[pattern]; exist {
		logger.L().Warn("event mux: multiple registrations for " + pattern)
		return
	}

	if mux.m == nil {
		mux.m = make(map[string]muxEntry)
	}

	e := muxEntry{
		h:       handler,
		pattern: pattern,
	}

	mux.m[pattern] = e
	mux.es = append(mux.es, e)

}

func (mux *TaskMux) HandleFunc(pattern string, handler HandlerFunc) {
	if handler == nil {
		return
	}
	mux.Handle(pattern, handler)
}

/////

func NotFound(ctx context.Context, task *Task) error {
	return fmt.Errorf("handler not found for task %q", task.Type())
}

func NotFoundHandler() Handler {
	return HandlerFunc(NotFound)
}
