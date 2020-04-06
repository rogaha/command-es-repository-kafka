package pkg

// TODO: create kafka implementation for this
type Provider interface {
	FetchAllEvents(batch int) (<-chan []Event, error)
	SendEvents(events []Event) error
}
