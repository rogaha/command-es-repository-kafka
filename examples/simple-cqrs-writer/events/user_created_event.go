package exampleevents

import (
	"encoding/json"
	"time"

	cerk "github.com/hetacode/command-es-repository-kafka"
)

type UserCreatedEvent struct {
	Type         string
	AggregatorId string
	CreateTime   string
	Version      int32

	Payload string

	FirstName string
	LastName  string
}

func (e *UserCreatedEvent) GetType() string {
	return e.Type
}
func (e *UserCreatedEvent) GetAggregatorId() string {
	return e.AggregatorId
}
func (e *UserCreatedEvent) GetCreateTime() time.Time {
	t, _ := time.Parse(time.RFC3339, e.CreateTime)

	return t
}
func (e *UserCreatedEvent) GetVersion() int32 {
	return e.Version
}
func (e *UserCreatedEvent) GetPayload() string {
	return e.Payload
}

func (e *UserCreatedEvent) LoadPayload() error {
	var jsonMap map[string]interface{}

	if err := json.Unmarshal([]byte(e.Payload), &jsonMap); err != nil {
		return err
	}

	e.Type = jsonMap["type"].(string)
	e.CreateTime = jsonMap["create_time"].(string)
	e.Version = int32(jsonMap["version"].(float64))
	e.FirstName = jsonMap["first_name"].(string)
	e.LastName = jsonMap["last_name"].(string)

	return nil
}

func (e *UserCreatedEvent) SavePayload() error {
	toJson := map[string]interface{}{
		"id":          e.AggregatorId,
		"type":        "UserCreatedEvent",
		"create_time": e.CreateTime,
		"version":     e.Version,
		"first_name":  e.FirstName,
		"last_name":   e.LastName,
	}
	bytesData, err := json.Marshal(toJson)

	if err != nil {
		return err
	}

	e.Payload = string(bytesData)

	return nil
}

func (e *UserCreatedEvent) InitBy(event cerk.Event) {
	e.Payload = event.GetPayload()
	e.AggregatorId = event.GetAggregatorId()
	e.Version = event.GetVersion()
}
