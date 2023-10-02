package models

import (
	"golang.org/x/crypto/bcrypt"
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
	if err := modelsApp.DB.GormDB.Where(condition).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// GetUserByID is a simple example of how get one user by id using gorm orm
func (u *User) GetUserByID(id int) (*User, error) {
	var user *User
	condition := User{ID: id}
	if err := modelsApp.DB.GormDB.Where(condition).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// InsertOneUser simply insert one user to the database
func (u *User) InsertOneUser(user *User) (int, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	user.Password = string(bytes)

	result := modelsApp.DB.GormDB.Create(user)
	if result.Error != nil {
		return 0, result.Error
	}

	return user.ID, nil
}

// InsertManyUsers insert many users to the database at the same time (create in batches)
func (u *User) InsertManyUsers(users []*User) (int64, []int, error) {
	var newUsersID []int

	for _, user := range users {
		bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return 0, nil, err
		}
		user.Password = string(bytes)
	}

	result := modelsApp.DB.GormDB.CreateInBatches(users, len(users))
	if result.Error != nil {
		return 0, nil, result.Error
	}

	for _, user := range users {
		newUsersID = append(newUsersID, user.ID)
	}

	return result.RowsAffected, newUsersID, nil
}
