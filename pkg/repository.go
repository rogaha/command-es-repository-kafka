package pkg

import "fmt"

// Repository it's an abstraction for database which keeps all entities (aggregators) in theirs last state
type Repository interface {
	SetProvider(provider Provider) // ??
	AddOrModifyEntity(entity Entity)
	GetEntity(id string) (Entity, error)
	Replay(events []Event) error
	GetUncommitedChanges() []Event
	Save(events []Event) error
}

// MemoryRepository is an basic implementation of Repository which keep data in the memory. This struct waiting for inheritance by own Repository.
// Inherited implementation should contains special methods for each needed case of event type
type MemoryRepository struct {
	entities       map[string]Entity
	eventsToCommit []Event
}

// AddOrModifyEntity just set new entity to collections of Entities.
// It will be replace if this id exists
func (r *MemoryRepository) AddOrModifyEntity(entity Entity) {
	r.entities[entity.GetId()] = entity
}

// GetEntity return current entity state provided by id
func (r *MemoryRepository) GetEntity(id string) (Entity, error) {
	entity := r.entities[id]
	if entity == nil {
		return nil, fmt.Errorf("Cannot find entity - id %s", id)
	}

	return entity, nil
}

// Replay method update state of entity/ies by provided events.
// This method should be override by true implementation with update cases for each event type
func (r *MemoryRepository) Replay(events []Event) error {
	return fmt.Errorf("Please implement this method in ith own way")
}

// GetUncommitedChanges get new events which were created by changes methods
func (r *MemoryRepository) GetUncommitedChanges() []Event {
	return r.eventsToCommit
}

// Save events - so to be honest, just send events to bus provider
func (r *MemoryRepository) Save(events []Event) error {
	// TODO: provider here
	r.eventsToCommit = make([]Event, 0)

	return nil
}

// NewMemoryRepository create empty initialized instance
func NewMemoryRepository() *MemoryRepository {
	repository := new(MemoryRepository)
	repository.entities = make(map[string]Entity)
	repository.eventsToCommit = make([]Event, 0)

	return repository
}
