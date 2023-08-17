package models

type BookInterface interface {
	GetAll() ([]*Book, error)
}

type UserInterface interface {
	GetOne() (*User, error)
}
