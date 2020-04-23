package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	cerk "github.com/hetacode/command-es-repository-kafka"
	examplecommandhandlers "github.com/hetacode/command-es-repository-kafka/examples/simple-cqrs-writer/command_handlers"
	examplecommands "github.com/hetacode/command-es-repository-kafka/examples/simple-cqrs-writer/commands"
	examplerepository "github.com/hetacode/command-es-repository-kafka/examples/simple-cqrs-writer/repository"
)

type MainContainer struct {
	repository *examplerepository.UsersRepository
}

func (h *MainContainer) handler(res http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Panic(err)
			res.WriteHeader(400)
		}
		log.Printf("Received commmand: %s", string(body))

		var jsonMap map[string]interface{}
		jsonData := body

		if err := json.Unmarshal(jsonData, &jsonMap); err != nil {
			log.Panic(err)
			res.WriteHeader(400)
		}

		// commands actions
		switch jsonMap["type"] {
		case "CreateUserCommand":
			var command *examplecommands.CreateUserCommand
			if err := json.Unmarshal(jsonData, &command); err != nil {
				log.Panic(err)
				res.WriteHeader(400)
			}

			if err := examplecommandhandlers.CreateUserCommandHandler(h.repository, command); err != nil {
				log.Panic(err)
				res.WriteHeader(400)
			}

		case "UpdateUserCommand":
			var command *examplecommands.UpdateUserCommand
			if err := json.Unmarshal(jsonData, &command); err != nil {
				log.Panic(err)
				res.WriteHeader(400)
			}

			if err := examplecommandhandlers.UpdateUserCommandHandler(h.repository, command); err != nil {
				log.Panic(err)
				res.WriteHeader(400)
			}
		}
	}
}

func main() {
	provider := cerk.NewKafkaProvider("example-topic", "example-app-group", "192.168.1.151:9092")
	defer provider.Close()

	repository := new(examplerepository.UsersRepository)
	repository.MemoryRepository = cerk.NewMemoryRepository()
	if err := repository.InitProvider(provider, repository); err != nil {
		panic(err)
	}

	h := &MainContainer{
		repository: repository,
	}

	http.HandleFunc("/", h.handler)
	http.ListenAndServe(":4000", nil)
}
