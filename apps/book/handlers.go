package book

import (
	"go-base-structure/pkg/json"
	"net/http"
)

// Home is a simple handler for book app
func Home(w http.ResponseWriter, r *http.Request) {
	users, err := bookApp.Models.Book.GetAll()
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		bookApp.Logger.Error("cannot get users from the database: ", err)
		return
	}

	err = json.WriteJSON(w, http.StatusOK, &users)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		bookApp.Logger.Error("unable to write json: ", err)
		return
	}
}
