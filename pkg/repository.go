package pkg

// Repository it's an abstraction for database which keeps all entities (aggregators) in theirs last state
type Repository interface {
	GetEntity(id string) (Entity, error)
	Replay(events []Event) error
	GetUncommitedChanges() []Event
	Save(event []Event) error
}
