package examplerepository

type UserEntity struct {
	ID        string
	FirstName string
	LastName  string
}

func (e *UserEntity) GetId() string {
	return e.ID
}
