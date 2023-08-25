package models

type TestUser struct {
	ID       int
	Email    string
	Password string
}

// GetUserByEmail for testing
func (u *TestUser) GetUserByEmail(email string) (*User, error) {
	var user *User
	return user, nil
}

// GetUserByID for testing
func (u *TestUser) GetUserByID(id int) (*User, error) {
	var user *User
	return user, nil
}
