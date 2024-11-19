package event

import (
	"reflect"
	"time"

	"github.com/google/uuid"
)

type EventType string

type EventStatus int

const (
	EventStatus_Idle EventStatus = iota
	EventStatus_OK
	EventStatus_Failed
)

//////////////////

type IEvent interface {
	AggregateID() string
	Version() int
	GetMetas() map[string]any
	SetMeta(string, any)
	EventType() string
	Event() any // payload
}

type RawEvent interface {
	GetEventType() EventType
	Summary() string
}

type Event struct {
	Version       int            `json:"version,omitempty"` // 事件版本
	EventId       string         `json:"event_id,omitempty"`
	EventType     EventType      `json:"event_type,omitempty"`
	AggregateId   string         `json:"aggregate_id,omitempty"`   // StreamID :聚合ID，GUID
	AggregateType string         `json:"aggregate_type,omitempty"` // StreamType
	Time          time.Time      `json:"time,omitempty"`
	Payload       any            `json:"payload,omitempty"`
	Metadata      *EventMetadata `json:"metadata,omitempty"`
}

// type EventMetadata struct {
// 	Uid         string `json:"uid,omitempty"` // Identity
// 	Username    string `json:"username,omitempty"`
// 	DeviceId    string `json:"device_id,omitempty"`
// 	DeviceModel string `json:"device_model,omitempty"`
// 	DeviceName  string `json:"device_name,omitempty"`
// 	DeviceType  string `json:"device_type,omitempty"`
// 	ViewType    string `json:"view_type,omitempty"`
// 	ViewId      string `json:"view_id,omitempty"`
// 	ViewName    string `json:"view_name,omitempty"`
// }

// func NewEvent(aggregateID string, eventType EventType, event any) Event {
func NewEvent(eventType EventType, event any) Event {

	return Event{
		EventId: uuid.NewString(),
		// AggregateId: aggregateID,
		// Version:     0,
		Time:      time.Now(),
		EventType: eventType,
		Payload:   event,
	}
}

func typeOf(i interface{}) string {
	return reflect.TypeOf(i).Elem().Name()
}
