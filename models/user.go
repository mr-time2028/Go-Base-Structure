package models

type User struct {
	ID       int
	Email    string
	Password string
}

func (u User) GetOne() (*User, error) {
	var user *User
	return user, nil
}
