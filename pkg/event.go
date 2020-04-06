package pkg

import "time"

// Event is an basic description for object event keep in event store and transfer between service (usually via some bus)
type Event interface {
	GetType() string
	GetAggregatorId() string
	GetCreateTime() time.Time
	GetVersion() int32

	LoadPayload() error
	SavePayload() error
}
