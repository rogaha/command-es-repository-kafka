package cerk

import (
	"encoding/json"
	"fmt"
	"time"
)

// Event is an basic description for object event keep in event store and transfer between service (usually via some bus)
type Event interface {
	GetType() string
	GetAggregatorId() string
	GetCreateTime() time.Time
	GetVersion() int32
	GetPayload() string

	LoadPayload() error
	SavePayload() error
	InitBy(event Event)
}

type GenericEvent struct {
	Type         string
	AggregatorId string
	CreateTime   string
	Version      int32

	Payload string
}

func (e *GenericEvent) GetType() string {
	return e.Type
}
func (e *GenericEvent) GetAggregatorId() string {
	return e.AggregatorId
}
func (e *GenericEvent) GetCreateTime() time.Time {
	t, _ := time.Parse(time.RFC3339, e.CreateTime)

	return t
}
func (e *GenericEvent) GetVersion() int32 {
	return e.Version
}
func (e *GenericEvent) GetPayload() string {
	return e.Payload
}

func (e *GenericEvent) LoadPayload() error {
	var jsonMap map[string]interface{}

	if err := json.Unmarshal([]byte(e.Payload), &jsonMap); err != nil {
		return err
	}

	e.Type = jsonMap["type"].(string)
	e.CreateTime = jsonMap["create_time"].(string)
	e.Version = int32(jsonMap["version"].(float64))

	return nil
}

func (e *GenericEvent) SavePayload() error {
	return fmt.Errorf("For GenericEvent implementation is unnecessary")
}

func (e *GenericEvent) InitBy(event Event) {
}
