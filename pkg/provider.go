package pkg

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// TODO: create kafka implementation for this
type Provider interface {
	FetchAllEvents(batch int) (<-chan []Event, error)
	SendEvents(events []Event) error
}

// KafkaProvider implemented provider for kafka
type KafkaProvider struct {
	topic     string
	servers   string
	groupName string
}

func (p *KafkaProvider) FetchAllEvents(batch int) (<-chan []Event, error) {
	rand.Seed(time.Now().UTC().UnixNano())

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":    p.servers,
		"group.id":             fmt.Sprintf("%s-replay-%d", p.groupName, rand.Int63n(10000)),
		"auto.offset.reset":    "earliest",
		"enable.partition.eof": true,
	})
	defer c.Close()

	if err != nil {
		return nil, err
	}

	c.SubscribeTopics([]string{p.topic}, nil)
	eventsChan := make(chan []Event)

	go func() {
		run := true
		currentMessageNo := 0
		events := make([]Event, 0)
		for run == true {
			ev := c.Poll(0)
			switch e := ev.(type) {
			case kafka.AssignedPartitions:
				fmt.Fprintf(os.Stderr, "%% %v\n", e)
				// c.Assign(e.Partitions)
			case kafka.RevokedPartitions:
				fmt.Fprintf(os.Stderr, "%% %v\n", e)
				// c.Unassign()
			case *kafka.Message:
				fmt.Printf("%% Message on %s:\n%s\n",
					e.TopicPartition, string(e.Value))

				// TODO: declare event

				if currentMessageNo >= batch {
					eventsChan <- events
					events = make([]Event, 0)
				}
				currentMessageNo++

			case kafka.PartitionEOF:
				fmt.Printf("%% Reached %v\n", e)
				run = false
			case kafka.Error:
				fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
				run = false
			}
		}

		if len(events) > 0 {
			eventsChan <- events
		}
	}()

	return eventsChan, nil
}

func (p *KafkaProvider) SendEvents(events []Event) error {
	return nil
}

// NewKafkaProvider create new instance of provider
func NewKafkaProvider(topic string, groupName string, servers string) Provider {
	provider := new(KafkaProvider)
	provider.servers = servers
	provider.topic = topic
	provider.groupName = groupName

	return provider
}
