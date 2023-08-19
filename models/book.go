package models

import (
	"context"
	"database/sql"
	"log"
	"time"
)

const dbTimeout = time.Second * 3

type Book struct {
	ID   int `gorm:"primaryKey"`
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
		return nil, err
	} else if err != nil {
		log.Println("Error query to plans table")
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var book Book
		err = rows.Scan(
			&book.ID,
			&book.Name,
		)
		if err != nil {
			log.Println("Error scanning a plan row")
			return nil, err
		}

		books = append(books, &book)
	}

	if err = rows.Err(); err != nil {
		log.Println("Error returned by plans rows")
		return nil, err
	}

	return books, nil
}
