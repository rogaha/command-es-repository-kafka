package test

import (
	"strings"
	"testing"
	"time"

	"github.com/hetacode/command-es-repository-kafka/pkg"
)

func Test_Replay_Method_For_MemoryRepository_Should_Be_Override(t *testing.T) {
	repository := new(pkg.MemoryRepository);
	
	err := repository.Replay(make([]pkg.Event, 0));
	if err == nil {
		t.Error("Error shouldn't be null")
	}

	if !strings.HasPrefix(err.Error(), "Please implement this") {
		t.Errorf("Wrong error - %s", err.Error())
	}  
}

// TODO: test LoadPayload, SavePayload of EntityCreatedEvent 

func Test_Replay_Function_In_MockRepository(t *testing.T) {
	event := &EntityCreatedEvent{
		AggregatorId: "1",
		Version: 1,
		Message: "It's just a test",
		CreateTime: time.Now().String(),
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