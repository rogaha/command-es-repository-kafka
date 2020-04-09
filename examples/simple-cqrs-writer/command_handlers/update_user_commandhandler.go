package examplecommandhandlers

import (
	cerk "github.com/hetacode/command-es-repository-kafka"
	examplecommands "github.com/hetacode/command-es-repository-kafka/examples/simple-cqrs-writer/commands"
	examplerepository "github.com/hetacode/command-es-repository-kafka/examples/simple-cqrs-writer/repository"
)

func UpdateUserCommandHandler(repository *examplerepository.UsersRepository, command *examplecommands.UpdateUserCommand) error {
	event, err := repository.Update(command.ID, command.FirstName, command.LastName)
	if err != nil {
		return err
	}
	if err := repository.Replay([]cerk.Event{event}); err != nil {
		return err
	}
	if err := repository.Save([]cerk.Event{event}); err != nil {
		return err
	}

	return nil
}
