package book

import (
	"net/http"
)

// Home is a simple handler for book app
func Home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("this is book app home page"))
}
