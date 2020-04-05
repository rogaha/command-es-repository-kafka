package test

import (
	"strings"
	"testing"
	"time"

	"github.com/hetacode/command-es-repository-kafka/pkg"
)

func Test_Replay_Method_For_MemoryRepository_Should_Be_Override(t *testing.T) {
	repository := new(pkg.MemoryRepository)

	err := repository.Replay(make([]pkg.Event, 0))
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
	repository.MemoryRepository = pkg.NewMemoryRepository()

	events := []pkg.Event{event}
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
