package tests

import "github.com/primus/primus/doma"

func init() {
	// Only the event for creating an invite has custom data.
	doma.RegisterEventData(CreateDemoEvent, func() doma.EventData {
		return &DemoCreateData{}
	})
}

const (
	CreateDemoEvent doma.EventType = "demo_created"
)

type DemoCreateData struct {
	Name string
}
