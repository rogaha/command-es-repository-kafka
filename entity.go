package cerk

// Entity is an interface represent one aggregate in repository
type Entity interface {
	GetId() string
}
