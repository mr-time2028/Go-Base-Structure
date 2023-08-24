package models

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

// dbTimeout is maximum of time that database operation can happen, use to create a context
const dbTimeout = time.Second * 3

// Book is a type for book table
type Book struct {
	ID   int `gorm:"primaryKey;autoIncrement"`
	Name string
}

// GetAll is an example of custom sql queries, you can also use of gorm's custom queries using GormDB.Raw()
func (b *Book) GetAll() ([]*Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `select * from books`

	var books []*Book

	rows, err := modelsApp.DB.SqlDB.QueryContext(ctx, query)
	if err == sql.ErrNoRows {
		return nil, errors.New("no rows found when query to books table. " + err.Error())
	} else if err != nil {
		return nil, errors.New("query to books table failed. " + err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var book Book
		err = rows.Scan(
			&book.ID,
			&book.Name,
		)
		if err != nil {
			return nil, errors.New("scanning book row failed. " + err.Error())
		}

		books = append(books, &book)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.New("something wrong with book rows. " + err.Error())
	}

	return books, nil
}
