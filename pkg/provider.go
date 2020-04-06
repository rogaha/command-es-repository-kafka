package pkg

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// Provider interface
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

// FetchAllEvents get all events from all partitions from specified topic
func (p *KafkaProvider) FetchAllEvents(batch int) (<-chan []Event, error) {
	rand.Seed(time.Now().UTC().UnixNano())

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":    p.servers,
		"group.id":             fmt.Sprintf("%s-replay-%d", p.groupName, rand.Int63n(10000)),
		"auto.offset.reset":    "earliest",
		"enable.partition.eof": true,
	})

	if err != nil {
		return nil, err
	}

	c.SubscribeTopics([]string{p.topic}, nil)
	eventsChan := make(chan []Event)

	// TODO: to handle fetching events from all partitions
	go func() {
		defer c.Close()

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
				// fmt.Printf("%% Message on %s:\n%s\n",
				// 	e.TopicPartition, string(e.Value))

				event := new(GenericEvent)
				event.AggregatorId = string(e.Key)
				event.Payload = string(e.Value)

				events = append(events, event)
				currentMessageNo++

				if currentMessageNo >= batch {
					eventsChan <- events
					events = make([]Event, 0)
					currentMessageNo = 0
				}

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

// SendEvents put messages on kafka topic
func (p *KafkaProvider) SendEvents(events []Event) error {
	pr, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": p.servers})

	if err != nil {
		return err
	}

	defer pr.Close()
	for _, e := range events {
		message := kafka.Message{
			Key:            []byte(e.GetAggregatorId()),
			TopicPartition: kafka.TopicPartition{Topic: &p.topic, Partition: kafka.PartitionAny},
			Value:          []byte(e.GetPayload()),
		}

		if err := pr.Produce(&message, nil); err != nil {
			return err
		}
	}

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
