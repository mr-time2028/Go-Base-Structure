package models

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
