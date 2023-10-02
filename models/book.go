package models

import (
	"context"
	"errors"
	"fmt"
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

// InsertOneBook simply insert one book to the database
func (b *Book) InsertOneBook(book *Book) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `insert into books (name) values ($1) returning id`

	tx, err := modelsApp.DB.SqlDB.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var newBookID int
	if err = tx.QueryRowContext(ctx, query, book.Name).Scan(&newBookID); err != nil {
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return newBookID, nil
}

// InsertManyBooks insert many books to the database at the same time (create in batches)
func (b *Book) InsertManyBooks(books []*Book) (int64, []int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `insert into books (name) values`

	var totalInserted int64
	for i, book := range books {
		query += fmt.Sprintf(` ('%s')`, book.Name)
		if i < len(books)-1 {
			query += ","
		}
		totalInserted += 1
	}
	query += ` returning id`

	tx, err := modelsApp.DB.SqlDB.BeginTx(ctx, nil)
	if err != nil {
		return 0, nil, err
	}
	defer tx.Rollback()

	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return 0, nil, err
	}
	defer rows.Close()

	var newBooksID []int
	for rows.Next() {
		var newBookID int
		if err = rows.Scan(&newBookID); err != nil {
			return 0, nil, err
		}
		newBooksID = append(newBooksID, newBookID)
	}

	if err = rows.Err(); err != nil {
		return 0, nil, err
	}

	if err = tx.Commit(); err != nil {
		return 0, nil, err
	}

	return totalInserted, newBooksID, nil
}
