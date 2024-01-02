package models

import (
	"testing"
)

func TestUser_GetUserByID(t *testing.T) {
	var testCases = []struct {
		name        string // name of the test
		userID      int    // specific user id that we want to get user with it from database
		expectedErr bool   // do we expect any error from this query to the database?
	}{
		{
			"get user by id",
			1,
			false,
		},
		{
			"no rows",
			3,
			true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// test function
			user := User{}
			_, err := user.GetUserByID(tc.userID)

			// validation
			if tc.expectedErr && err == nil || !tc.expectedErr && err != nil {
				t.Errorf("unexpected error: expectedErr is %t, err is %s", tc.expectedErr, err.Error())
			}
		})
	}
}

func TestUser_GetUserByEmail(t *testing.T) {
	var testCases = []struct {
		name        string // name of the test
		userEmail   string // specific user email that we want to get user with it from database
		expectedErr bool   // do we expect any error from this query to the database?
	}{
		{
			"get user by email",
			"John@test.com",
			false,
		},
		{
			"no rows",
			"Benjamin@test.com",
			true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// test function
			user := User{}
			_, err := user.GetUserByEmail(tc.userEmail)

			// validation
			if tc.expectedErr && err == nil || !tc.expectedErr && err != nil {
				t.Errorf("unexpected error: expectedErr is %v, err is %v", tc.expectedErr, err)
			}
		})
	}
}

func TestUser_InsertOneUser(t *testing.T) {
	var testCases = []struct {
		name        string // name of the test
		user        *User  // user we want to insert to the database
		expectedErr bool   // do we expect any error from this query to the database?
	}{
		{
			"insert one user",
			&User{ID: 3, Email: "Benjamin@test.com", FirstName: "Benjamin", LastName: "Smith", Password: "davidPass"},
			false,
		},
		{
			"insert one user (duplicate id)",
			&User{ID: 1, Email: "Benjamin@test.com", FirstName: "Benjamin", LastName: "Smith", Password: "davidPass"},
			true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// test function
			user := User{}
			userID, err := user.InsertOneUser(tc.user)

			// validation
			if tc.expectedErr && err == nil || !tc.expectedErr && err != nil {
				t.Errorf("unexpected error: expectedErr is %v, err is %v", tc.expectedErr, err)
			} else if err == nil && userID == 0 {
				t.Errorf("expected user id grather than 0 but it is not")
			}

			err = resetTestDB()
			if err != nil {
				logr.Fatal("failed to reset the database")
			}
		})
	}
}

func TestUser_InsertManyUsers(t *testing.T) {
	var testCases = []struct {
		name        string  // name of the test
		users       []*User // users we want to insert to the database
		expectedErr bool    // do we expect any error from this query to the database?
	}{
		{
			"insert many users",
			[]*User{
				{ID: 3, Email: "Mary@test.com", FirstName: "Mary", LastName: "Jane", Password: "MaryPass"},
				{ID: 4, Email: "John@test.com", FirstName: "John", LastName: "Montgomery", Password: "JohnPass"},
			},
			false,
		},
		{
			"insert many users (duplicate id)",
			[]*User{
				{ID: 1, Email: "Mary@test.com", FirstName: "Mary", LastName: "Jane", Password: "MaryPass"},
				{ID: 4, Email: "John@test.com", FirstName: "John", LastName: "Montgomery", Password: "JohnPass"},
			},
			true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// test function
			user := User{}
			rowAffected, usersID, err := user.InsertManyUsers(tc.users)

			// validation
			if tc.expectedErr && err == nil || !tc.expectedErr && err != nil {
				t.Errorf("unexpected error: expectedErr is %v, err is %v", tc.expectedErr, err)
			} else if err == nil && rowAffected != int64(len(tc.users)) {
				t.Errorf("expected rowAffected equal to %d, but it is equal to: %d", int64(len(tc.users)), rowAffected)
			} else if err == nil && rowAffected != int64(len(usersID)) {
				t.Errorf("expected rowAffected equals to usersID length but it is not")
			}

			err = resetTestDB()
			if err != nil {
				logr.Fatal("failed to reset the database")
			}
		})
	}
}

func TestUser_CheckIfExistsUser(t *testing.T) {
	var testCases = []struct {
		name  string
		user  *User
		expectedErr bool
		isExists bool
	} {
		{
			"user exists",
			&User{Email: "John@test.com"},
			false,
			true,
		},
		{
			"user does not exists",
			&User{Email: "noexists@test.com"},
			false,
			false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			user := User{}
			isExists, err := user.CheckIfExistsUser(tc.user.Email)

			// validation
			if tc.expectedErr && err == nil || !tc.expectedErr && err != nil {
				t.Errorf("unexpected error: expectedErr is %v, err is %v", tc.expectedErr, err)
			} else if isExists != tc.isExists {
				t.Errorf("expected user exists %v but got user exists %v", tc.isExists, isExists)
			}
		})
	}
}
