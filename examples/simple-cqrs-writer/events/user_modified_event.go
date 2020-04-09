package exampleevents

import (
	"encoding/json"
	"time"

	"github.com/hetacode/command-es-repository-kafka/pkg"
)

type UserModifiedEvent struct {
	Type         string
	AggregatorId string
	CreateTime   string
	Version      int32

	Payload string

	FirstName string
	LastName  string
}

func (e *UserModifiedEvent) GetType() string {
	return e.Type
}
func (e *UserModifiedEvent) GetAggregatorId() string {
	return e.AggregatorId
}
func (e *UserModifiedEvent) GetCreateTime() time.Time {
	t, _ := time.Parse(time.RFC3339, e.CreateTime)

	return t
}
func (e *UserModifiedEvent) GetVersion() int32 {
	return e.Version
}
func (e *UserModifiedEvent) GetPayload() string {
	return e.Payload
}

func (e *UserModifiedEvent) LoadPayload() error {
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

func (e *UserModifiedEvent) SavePayload() error {
	toJson := map[string]interface{}{
		"id":          e.AggregatorId,
		"type":        "UserModifiedEvent",
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

func (e *UserModifiedEvent) InitBy(event pkg.Event) {
	e.Payload = event.GetPayload()
	e.AggregatorId = event.GetAggregatorId()
	e.Version = event.GetVersion()
}
