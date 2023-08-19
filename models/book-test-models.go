package models

type TestBook struct {
	ID   int
	Name string
}

// GetAll for testing
func (b *TestBook) GetAll() ([]*Book, error) {
	var books []*Book
	return books, nil
}
