package models

import (
	"errors"
	"gorm.io/gorm"
)

type User struct {
	ID        int
	Email     string
	FirstName string
	LastName  string
	Password  string
}

// GetUserByEmail is a simple example of how get one user by email using gorm orm
func (u *User) GetUserByEmail(email string) (*User, error) {
	var user *User
	condition := User{Email: email}
	err := modelsApp.DB.GormDB.Where(condition).First(&user).Error

	if err == gorm.ErrRecordNotFound {
		return nil, errors.New("user with this email not found " + err.Error())
	} else if err != nil {
		return nil, err
	} else {
		return user, nil
	}
}

// GetUserByID is a simple example of how get one user by id using gorm orm
func (u *User) GetUserByID(id int) (*User, error) {
	var user *User
	condition := User{ID: id}
	err := modelsApp.DB.GormDB.Where(condition).First(&user).Error

	if err == gorm.ErrRecordNotFound {
		return nil, errors.New("user with this id not found " + err.Error())
	} else if err != nil {
		return nil, err
	} else {
		return user, nil
	}
}
