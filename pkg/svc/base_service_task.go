package svc

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/nats-io/nats.go"
	"github.com/primus/primus/core/common/task"
	"github.com/primus/primus/pkg/eventbus"
	"github.com/primus/primus/pkg/logger"
)

// StartTask implements Service.
func (s *BaseService) StartTask(target string, task task.Task) error {
	if target == "" {
		return fmt.Errorf("target null")
	}

	return s.bus.Publish(fmt.Sprintf("task.%s", target), task)
}

// SubscribeTask implements Service.
func (s *BaseService) SubscribeTask(fn TaskHandler) error {

	if s.bus == nil {
		return nil
	}
	if fn == nil {
		return nil
	}

	s.bus.Subscribe(fmt.Sprintf("task.%s", s.name), s.handleTask(fn))
	return nil

}

func (s *BaseService) handleTask(fn TaskHandler) func(msg *nats.Msg) {
	return func(msg *nats.Msg) {
		taskx := task.Task{}
		err := json.Unmarshal(msg.Data, &taskx)
		if err != nil {
			logger.L().Error("----service:", s.name, ", handle task, unmarshal err:", err)
			return
		}

		if _, exists := s.tasks.LoadOrStore(taskx.TaskId, taskx); exists {
			logger.L().Warnf("service [%s] task %s is running alreay , task type %s", s.name, taskx.TaskId, taskx.TaskType)
			return
		}

		go func() {
			defer func() {
				s.tasks.Delete(taskx.TaskId)
			}()
			logger.L().Info("----service:", s.name, ", receive task bus msg: ", taskx.TaskType)

			taskCtx := task.WrapTaskContext(context.Background(), taskx.TaskId, taskx.TaskType, nil)

			res, err := s.performTask(taskCtx, taskx, fn)
			if err != nil {
				logger.L().Error("----service:", s.name, ", handle task ", taskx.TaskType, " err:", err)
				s.bus.Publish(msg.Reply, err.Error())
				return
			}
			bb, _ := json.Marshal(res)
			s.bus.Publish(msg.Reply, string(bb))
		}()

	}
}

func (s *BaseService) handleFailedMessage(ctx context.Context, taskId string, err error) error {
	// return s.bus.Publish()
	return nil
}

type taskProgress struct {
	bus eventbus.EventBus
}

// PublishProgress implements task.TaskProgress.
func (t *taskProgress) PublishProgress(success int32, total int32) error {
	if t.bus == nil {
		return nil
	}
	bb, _ := json.Marshal(struct {
		Success int32 `json:"success"`
		Total   int32 `json:"totoal"`
	}{})
	return t.bus.Publish("reply+progress", bb)
}

func newTaskProgress() task.TaskProgress {
	return &taskProgress{}
}

func (s *BaseService) performTask(ctx context.Context, task task.Task, fn TaskHandler) (res any, err error) {
	defer func() {
		if x := recover(); x != nil {
			errMsg := string(debug.Stack())

			logger.L().Errorf("recovering from panic. See the stack trace below for details:\n%s", errMsg)
			_, file, line, ok := runtime.Caller(1) // skip the first frame (panic itself)
			if ok && strings.Contains(file, "runtime/") {
				// The panic came from the runtime, most likely due to incorrect
				// map/slice usage. The parent frame should have the real trigger.
				_, file, line, ok = runtime.Caller(2)
			}

			// Include the file and line number info in the error, if runtime.Caller returned ok.
			if ok {
				err = fmt.Errorf("panic [%s:%d]: %v", file, line, x)
			} else {
				err = fmt.Errorf("panic: %v", x)
			}
		}
	}()
	return fn(ctx, task)
}
