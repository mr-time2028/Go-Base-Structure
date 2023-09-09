package models

import (
	"gorm.io/gorm"
)

type TestUser struct {
	ID       int
	Email    string
	Password string
}

// GetUserByEmail for testing
func (u *TestUser) GetUserByEmail(email string) (*User, error) {
	if email == "norows@test.com" {
		return nil, gorm.ErrRecordNotFound
	}
	user := &User{
		ID:        1,
		FirstName: "somename",
		LastName:  "somelast",
		Password:  "$2a$10$G54ZltuaC.70vDH7f831FeNHmwe0FnVY82M9RxoUsJOVjOwRDn.tS", // password is "testPass"
	}
	return user, nil
}

// GetUserByID for testing
func (u *TestUser) GetUserByID(id int) (*User, error) {
	user := &User{
		ID:        1,
		FirstName: "testname",
		LastName:  "testlast",
		Password:  "$2a$10$G54ZltuaC.70vDH7f831FeNHmwe0FnVY82M9RxoUsJOVjOwRDn.tS", // password is "testPass"
	}
	return user, nil
}
