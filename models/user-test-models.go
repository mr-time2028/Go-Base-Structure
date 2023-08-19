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
