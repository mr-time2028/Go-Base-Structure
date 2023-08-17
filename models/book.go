package models

import (
	"context"
	"database/sql"
	"log"
	"time"
)

type Book struct {
	ID   int    `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
}

// GetAll is an example of custom queries
func (b *Book) GetAll() ([]*Book, error) {
	//return nil, nil
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	query := `select * from books`

	var books []*Book

	rows, err := SqlDB.QueryContext(ctx, query)
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
