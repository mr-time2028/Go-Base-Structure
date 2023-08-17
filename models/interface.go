package models

type BookInterface interface {
	GetAll() ([]*Book, error)
}

type UserInterface interface {
	GetOne(email string) (*User, error)
}
