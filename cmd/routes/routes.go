package routes

import (
	"github.com/go-chi/chi/v5"
	bookRoute "go-base-structure/apps/book"
	userRoute "go-base-structure/apps/user"
	"net/http"
)

// Routes aggregate all routes from all apps
func Routes() http.Handler {
	mux := chi.NewRouter()

	mux.Mount("/books", bookRoute.Routes())
	mux.Mount("/users", userRoute.Routes())

	return mux
}
