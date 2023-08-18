package models

type TestUser struct {
	ID       int
	Email    string
	Password string
}

// GetOne for testing
func (u *TestUser) GetOne(email string) (*User, error) {
	var user *User
	return user, nil
}

type TestBook struct {
	ID   int `gorm:"primaryKey"`
	Name string
}

// GetAll for testing
func (b *TestBook) GetAll() ([]*Book, error) {
	var books []*Book
	return books, nil
}
