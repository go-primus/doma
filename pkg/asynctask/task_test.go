package asynctask

import (
	"context"
	"fmt"
	"testing"
)

func TestTask(t *testing.T) {

	task := NewAsyncTask(context.Background(), "task-demo", func(ctx context.Context) error {
		fmt.Println("hello")
		return nil
	})
	task.Execute()

}
