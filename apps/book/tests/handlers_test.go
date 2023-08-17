package tests

import (
	"go-base-structure/apps/book"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHome(t *testing.T) {
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/books/book-home", nil)
	handler := http.HandlerFunc(book.Home)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Error("there is an error")
	}
}
