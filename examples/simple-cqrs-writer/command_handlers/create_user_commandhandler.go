package examplecommandhandlers

import (
	examplecommands "github.com/hetacode/command-es-repository-kafka/examples/simple-cqrs-writer/commands"
	examplerepository "github.com/hetacode/command-es-repository-kafka/examples/simple-cqrs-writer/repository"
	"github.com/hetacode/command-es-repository-kafka/pkg"
)

func CreateUserCommandHandler(repository *examplerepository.UsersRepository, command *examplecommands.CreateUserCommand) error {
	event, err := repository.Create(command.FirstName, command.LastName)
	if err != nil {
		return err
	}
	if err := repository.Replay([]pkg.Event{event}); err != nil {
		return err
	}
	if err := repository.Save([]pkg.Event{event}); err != nil {
		return err
	}

	return nil
}
