package models

import (
	"errors"
	"gorm.io/gorm"
)

type User struct {
	ID       int
	Email    string
	Password string
}

// GetOne is a simple example of how get one user using gorm orm
func (u *User) GetOne(email string) (*User, error) {
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
