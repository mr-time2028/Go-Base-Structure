package models

import (
	"context"
	"errors"
	"time"
)

// dbTimeout is maximum of time that database operation can happen, use to create a context
const dbTimeout = time.Second * 3

// Book is a type for book table
type Book struct {
	ID   int    `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"unique"`
}

// GetAll is an example of custom sql queries, you can also use of gorm's custom queries using GormDB.Raw()
func (b *Book) GetAll() ([]*Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select * from books`

	var books []*Book

	rows, err := modelsApp.DB.SqlDB.QueryContext(ctx, query)
	if err != nil {
		return nil, errors.New("failed query to books table: " + err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var book Book
		err = rows.Scan(
			&book.ID,
			&book.Name,
		)
		if err != nil {
			return nil, errors.New("failed scanning book row: " + err.Error())
		}

		books = append(books, &book)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.New("something wrong with book rows: " + err.Error())
	}

	return books, nil
}
