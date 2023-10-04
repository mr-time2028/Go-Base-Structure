package book

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestHome simply test home handler
func TestHome(t *testing.T) {
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/books/book-home", nil)
	handler := http.HandlerFunc(Home)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, but got %d", http.StatusOK, rr.Code)
	}

	if rr.Body == nil {
		t.Error("unexpected body, body is nil")
	}
}
