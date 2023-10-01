package book

import (
	"go-base-structure/pkg/json"
	"net/http"
)

// Home is a simple handler for book app
func Home(w http.ResponseWriter, r *http.Request) {
	books, err := bookApp.Models.Book.GetAll()
	if err != nil {
		bookApp.Logger.ServerError(w, "cannot get books from the database", err)
		return
	}

	if err = json.WriteJSON(w, http.StatusOK, &books); err != nil {
		bookApp.Logger.ServerError(w, "unable to write json", err)
		return
	}
}
