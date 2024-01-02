package user

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

// Routes contains all routes of user app
func Routes() http.Handler {
	mux := chi.NewRouter()

	mux.Post("/login", Login)
	mux.Post("/register", Register)
	mux.Post("/refresh", RefreshToken)

	return mux
}
