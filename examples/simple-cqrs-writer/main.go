package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	examplecommands "github.com/hetacode/command-es-repository-kafka/examples/simple-cqrs-writer/commands"
)

type CommandsHandler struct {
	repository *UsersRepository
}

func (h *CommandsHandler) handler(res http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Panic(err)
			res.WriteHeader(400)
		}

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
			// TODO: command handler action

		case "UpdateUserCommand":
			var command *examplecommands.UpdateUserCommand
			if err := json.Unmarshal(jsonData, &command); err != nil {
				log.Panic(err)
				res.WriteHeader(400)
			}

			// TODO: command handler action
		}
	}
}

func main() {
	h := &CommandsHandler{
		repository: new(UsersRepository),
	}

	http.HandleFunc("/", h.handler)
}
