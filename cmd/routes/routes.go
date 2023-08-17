package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	bookRoute "go-base-structure/apps/book"
	userRoute "go-base-structure/apps/user"
	"go-base-structure/cmd/middlewares"
	"net/http"
)

// Routes aggregate all routes from all apps
func Routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(middlewares.CORSMiddleware)

	mux.Mount("/books", bookRoute.Routes())
	mux.Mount("/users", userRoute.Routes())

	return mux
}
