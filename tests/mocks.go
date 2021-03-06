package test

import (
	"encoding/json"
	"time"

	cerk "github.com/hetacode/command-es-repository-kafka"
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
func (e *EntityCreatedEvent) GetPayload() string {
	return e.Payload
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

func (e *EntityCreatedEvent) InitBy(event cerk.Event) {
	e.Payload = event.GetPayload()
	e.AggregatorId = event.GetAggregatorId()
	e.Version = event.GetVersion()
}

type MockEntity struct {
	Id string

	Message string
}

func (e *MockEntity) GetId() string {
	return e.Id
}

type MockProvider struct {
	Events     []cerk.Event
	initEvents []cerk.Event
}

func (p *MockProvider) SetInitEvents(events []cerk.Event) {
	p.initEvents = events
}

func (p *MockProvider) FetchAllEvents(batch int) (<-chan []cerk.Event, error) {
	c := make(chan []cerk.Event)
	go func() {
		c <- p.initEvents
		close(c)
	}()
	return c, nil
}

func (p *MockProvider) SendEvents(events []cerk.Event) error {
	p.Events = events
	return nil
}

func (p *MockProvider) Close() {

}

type MockRepository struct {
	*cerk.MemoryRepository
}

func (r *MockRepository) Replay(events []cerk.Event) error {
	for _, e := range events {
		e.LoadPayload()
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

func (r *MockRepository) CreateFakeEvent() {
	event := &EntityCreatedEvent{
		AggregatorId: "1",
		Version:      1,
		Payload:      `{"message":"fake", "createTime":"2009-11-10T23:00:00Z"}`,
	}
	r.AddNewEvent(event)
}
