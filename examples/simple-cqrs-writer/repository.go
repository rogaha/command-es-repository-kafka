package main

import (
	exampleevents "github.com/hetacode/command-es-repository-kafka/examples/simple-cqrs-writer/events"
	"github.com/hetacode/command-es-repository-kafka/pkg"
)

type UsersRepository struct {
	*pkg.MemoryRepository
}

func (r *UsersRepository) Replay(events []pkg.Event) error {
	for _, e := range events {
		e.LoadPayload()
		switch e.GetType() {
		case "UserCreatedEvent":
			event := e.(*exampleevents.UserCreatedEvent)
			entity := new(UserEntity)
			entity.ID = e.GetAggregatorId()
			entity.FirstName = event.FirstName
			entity.LastName = event.LastName
			r.AddOrModifyEntity(entity)
		case "UserModifiedEvent":
			event := e.(*exampleevents.UserModifiedEvent)
			entity, err := r.GetEntity(event.GetAggregatorId())
			if err != nil {
				return err
			}
			userEntity := entity.(*UserEntity)
			userEntity.FirstName = event.FirstName
			userEntity.LastName = event.LastName
			r.AddOrModifyEntity(userEntity)
		}
	}

	return nil
}
