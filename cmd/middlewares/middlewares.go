package middlewares

import (
	"github.com/go-chi/cors"
	"net/http"
)

// CORSMiddleware handle CORS stuff
func CORSMiddleware(next http.Handler) http.Handler {
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	return corsMiddleware.Handler(next)
}
