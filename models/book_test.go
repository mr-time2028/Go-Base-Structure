package models

import (
	"testing"
)

func TestBook_GetAll(t *testing.T) {
	var testCases = []struct {
		name        string // name of the test
		expectedErr bool   // do we expect any error from this query to the database?
	}{
		{
			"get all books",
			false,
		},
		{
			"no rows",
			false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// test function
			book := Book{}
			result, err := book.GetAll()

			// validation
			if tc.expectedErr && err == nil || !tc.expectedErr && err != nil {
				t.Errorf("unexpected error: expectedErr is %v, err is %v", tc.expectedErr, err)
			} else if len(result) != 2 {
				t.Errorf("expected two book record(s) in the database but it is %d record(s)", len(result))
			}
		})
	}
}

func TestBook_InsertOneBook(t *testing.T) {
	var testCases = []struct {
		name        string // name of the test
		book        *Book  // book we want to insert to the database
		expectedErr bool   // do we expect any error from this query to the database?
	}{
		{
			"insert one book",
			&Book{Name: "The Lord of the Rings"},
			false,
		},
		{
			"insert one book (duplicate id)",
			&Book{ID: 1, Name: "A Song of Ice and Fire"},
			false,
		},
		{
			"insert one book (duplicate name)",
			&Book{Name: "Harry Potter"},
			true, // because we use "name" has a unique gorm tag
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// test function
			book := Book{}
			newBookID, err := book.InsertOneBook(tc.book)

			// validation
			if tc.expectedErr && err == nil || !tc.expectedErr && err != nil {
				t.Errorf("unexpected error: expectedErr is %v, err is %v", tc.expectedErr, err)
			} else if err == nil && newBookID == 0 {
				t.Errorf("expected new book id grather than 0 but it is not")
			}

			err = resetTestDB()
			if err != nil {
				logr.Fatal("failed to reset the database")
			}
		})
	}
}

func TestBook_InsertManyBooks(t *testing.T) {
	var testCases = []struct {
		name        string  // name of the test
		books       []*Book // books we want to insert to the database
		expectedErr bool    // do we expect any error from this query to the database?
	}{
		{
			"insert many books",
			[]*Book{
				{Name: "The Lord of the Rings"},
				{Name: "A Song of Ice and Fire"},
			},
			false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// test function
			book := Book{}
			rowAffected, booksID, err := book.InsertManyBooks(tc.books)

			// validation
			if tc.expectedErr && err == nil || !tc.expectedErr && err != nil {
				t.Errorf("unexpected error: expectedErr is %v, err is %v", tc.expectedErr, err)
			} else if err == nil && rowAffected != int64(len(tc.books)) {
				t.Errorf("expected rowAffected equal to %d, but it is equal to: %d", int64(len(tc.books)), rowAffected)
			} else if err == nil && rowAffected != int64(len(booksID)) {
				t.Errorf("expected rowAffected equals to booksID length but it is not")
			}

			err = resetTestDB()
			if err != nil {
				logr.Fatal("failed to reset the database")
			}
		})
	}
}
