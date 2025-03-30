package tests

import "github.com/go-primus/doma/runtime/core"

func init() {
	// Only the event for creating an invite has custom data.
	core.RegisterEventData(CreateDemoEvent, func() core.EventData {
		return &DemoCreateData{}
	})
}

const (
	CreateDemoEvent core.EventType = "demo_created"
)

type DemoCreateData struct {
	Name string
}
