package test

import (
	"encoding/json"
	"time"

	"github.com/hetacode/command-es-repository-kafka/pkg"
)

type EntityCreatedEvent struct {
	AggregatorId string
	CreateTime   string
	Version      int32
	Payload      string

	Message string
}

func (e *EntityCreatedEvent) GetType() string {
	return "EntityCreatedEvent"
}
func (e *EntityCreatedEvent) GetAggregatorId() string {
	return e.AggregatorId
}
func (e *EntityCreatedEvent) GetCreateTime() time.Time {
	value, _ := time.Parse(time.RFC3339, e.CreateTime)
	return value
}
func (e *EntityCreatedEvent) GetVersion() int32 {
	return e.Version
}
func (e *EntityCreatedEvent) LoadPayload() error {
	var jsonMap map[string]interface{}
	bytesData := []byte(e.Payload)
	if err := json.Unmarshal(bytesData, &jsonMap); err != nil {
		return err
	}
	e.CreateTime = jsonMap["createTime"].(string)
	e.Message = jsonMap["message"].(string)

	return nil
}
func (e *EntityCreatedEvent) SavePayload() error {
	toJson := map[string]interface{}{
		"createTime": e.CreateTime,
		"message":    e.Message,
	}
	bytesData, err := json.Marshal(toJson)
	if err != nil {
		return err
	}
	e.Payload = string(bytesData)
	return nil
}

type MockEntity struct {
	Id string

	Message string
}

func (e *MockEntity) GetId() string {
	return e.Id
}

type MockRepository struct {
	*pkg.MemoryRepository
}

func (r *MockRepository) Replay(events []pkg.Event) error {
	for _, e := range events {
		switch e.GetType() {
		case "EntityCreatedEvent":
			entity := new(MockEntity)
			entity.Id = e.GetAggregatorId()
			entity.Message = e.(*EntityCreatedEvent).Message
			r.AddOrModifyEntity(entity)
		}
	}

	return nil
}
