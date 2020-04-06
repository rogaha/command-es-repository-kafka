package main

import (
	"flag"
	"fmt"

	"github.com/hetacode/command-es-repository-kafka/pkg"
)

func fetch(server string, topic string, group string) {
	provider := pkg.NewKafkaProvider(topic, group, server)
	eventsBatch, err := provider.FetchAllEvents(10)

	if err != nil {
		panic(err)
	}

	for events := range eventsBatch {
		for _, event := range events {
			fmt.Printf("%s\n", event.GetPayload())
		}
	}
}

func fill(server string, topic string, group string) {

}

func main() {
	action := flag.String("action", "fill", "fill | fetch events")
	topic := flag.String("topic", "test", "a kafka topic")
	server := flag.String("server", "localhost:9092", "a kafka server")
	group := flag.String("group", "test-group", "a kafka consumer group name")

	flag.Parse()

	switch *action {
	case "fetch":
		fetch(*server, *topic, *group)
	case "fill":
		fill(*server, *topic, *group)
	}

}
