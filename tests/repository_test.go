package test

import (
	"strings"
	"testing"
	"time"

	cerk "github.com/hetacode/command-es-repository-kafka"
)

func Test_Replay_Method_For_MemoryRepository_Should_Be_Override(t *testing.T) {
	repository := new(cerk.MemoryRepository)

	err := repository.Replay(make([]cerk.Event, 0))
	if err == nil {
		t.Error("Error shouldn't be null")
	}

	if !strings.HasPrefix(err.Error(), "Please implement this") {
		t.Errorf("Wrong error - %s", err.Error())
	}
}

func Test_Event_LoadPayload(t *testing.T) {
	event := &EntityCreatedEvent{
		AggregatorId: "1",
		Version:      1,
		Payload:      `{"message":"Just test", "createTime":"2009-11-10T23:00:00Z"}`,
	}
	if err := event.LoadPayload(); err != nil {
		t.Fatal(err.Error())
	}

	if event.Message != "Just test" {
		t.Fatalf("Wrong Message %s", event.Message)
	}

	if event.GetCreateTime().Year() != 2009 {
		t.Fatalf("Wrong CreateTime %s", event.GetCreateTime().Local().String())
	}
}

func Test_Event_SavePayload(t *testing.T) {
	event := &EntityCreatedEvent{
		AggregatorId: "1",
		Version:      1,
		Message:      "It's just a test",
		CreateTime:   time.Now().String(),
	}
	if err := event.SavePayload(); err != nil {
		t.Fatal(err.Error())
	}

	if len(event.Payload) == 0 {
		t.Fatal("Payload is empty")
	}

	if !strings.HasPrefix(event.Payload, `{"createTime":`) {
		t.Fatalf("Wrong payload - %s", event.Payload)
	}
}

func Test_Replay_Function_In_MockRepository(t *testing.T) {
	event := &EntityCreatedEvent{
		AggregatorId: "1",
		Version:      1,
		Message:      "It's just a test",
		CreateTime:   time.Now().String(),
	}

	repository := new(MockRepository)
	repository.MemoryRepository = cerk.NewMemoryRepository()

	events := []cerk.Event{event}
	if err := repository.Replay(events); err != nil {
		t.Fatal(err.Error())
	}

	entity, err := repository.GetEntity("1")
	if err != nil {
		t.Fatal(err.Error())
	}

	if entity == nil {
		t.Fatal("Cannot find entity")
	}

	if !strings.HasPrefix(entity.(*MockEntity).Message, "It's just") {
		t.Fatalf("Wrong Entity body - %s", entity.(*MockEntity).Message)
	}
}

func Test_InitProvider_And_Check_Generated_Entity(t *testing.T) {
	event := &EntityCreatedEvent{
		AggregatorId: "1",
		Version:      1,
		Payload:      `{"message":"Just test", "createTime":"2009-11-10T23:00:00Z"}`,
	}
	provider := new(MockProvider)
	provider.SetInitEvents([]cerk.Event{event})

	repository := new(MockRepository)
	repository.MemoryRepository = cerk.NewMemoryRepository()

	if err := repository.InitProvider(provider, repository); err != nil {
		t.Fatal(err.Error())
	}

	entity, err := repository.GetEntity("1")
	if err != nil {
		t.Fatal(err.Error())
	}

	if entity == nil {
		t.Fatal("Cannot find Entity")
	}

	if entity.(*MockEntity).Message != "Just test" {
		t.Fatalf("Entity has incorrect message - %s", entity.(*MockEntity).Message)
	}
}

func Test_Save_Events(t *testing.T) {
	event := &EntityCreatedEvent{
		AggregatorId: "1",
		Version:      1,
		Payload:      `{"message":"Just test", "createTime":"2009-11-10T23:00:00Z"}`,
	}

	provider := new(MockProvider)

	repository := new(MockRepository)
	repository.MemoryRepository = cerk.NewMemoryRepository()

	if err := repository.InitProvider(provider, repository); err != nil {
		t.Fatal(err.Error())
	}

	if err := repository.Save([]cerk.Event{event}); err != nil {
		t.Fatal(err.Error())
	}

	if len(provider.Events) == 0 {
		t.Fatal("Provider events shouldn't be empty")
	}
}

func Test_Init_Event_By_Other_Event(t *testing.T) {
	event := &EntityCreatedEvent{
		AggregatorId: "1",
		Version:      1,
		Payload:      `{"message":"Just test", "createTime":"2009-11-10T23:00:00Z"}`,
	}
	if err := event.LoadPayload(); err != nil {
		t.Fatal(err)
	}

	eventToCompare := new(EntityCreatedEvent)
	eventToCompare.InitBy(event)
	if err := eventToCompare.LoadPayload(); err != nil {
		t.Fatal(err)
	}

	if event.GetAggregatorId() != eventToCompare.GetAggregatorId() {
		t.Fatal("Events have different ids")
	}
	if event.Message != eventToCompare.Message {
		t.Fatal("Events have different messages")
	}
}
