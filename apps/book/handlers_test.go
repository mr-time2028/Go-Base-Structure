package book

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHome(t *testing.T) {
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/books/book-home", nil)
	handler := http.HandlerFunc(Home)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Error("there is an error")
	}
}
