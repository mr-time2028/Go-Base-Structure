package book

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

// Routes contains all routes of book app
func Routes() http.Handler {
	mux := chi.NewRouter()

	mux.Get("/book-home", Home)

	return mux
}
