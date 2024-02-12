package models

type BookInterface interface {
	GetAll() ([]*Book, error)
	InsertOneBook(*Book) (int, error)
	InsertManyBooks([]*Book) (int64, []int, error)
}

type UserInterface interface {
	CheckIfExistsUser(email string) (bool, error)
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int) (*User, error)
	InsertOneUser(user *User) (int, error)
	InsertManyUsers(user []*User) (int64, []int, error)
}

type ModelManager struct {
	Book BookInterface
	User UserInterface
}

func NewModels() *ModelManager {
	return &ModelManager{
		Book: &Book{},
		User: &User{},
	}
}
