package models

// TestBook for tests
type TestBook struct {
	ID   int
	Name string
}

// GetAll for tests
func (b *TestBook) GetAll() ([]*Book, error) {
	var books []*Book
	return books, nil
}
