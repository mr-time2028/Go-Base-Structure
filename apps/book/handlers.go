package book

import (
	"go-base-structure/pkg/json"
	"net/http"
)

// Home is a simple handler for book app
func Home(w http.ResponseWriter, r *http.Request) {
	users, err := bookApp.Models.Book.GetAll()
	if err != nil {
		bookApp.Logger.ErrorLog.Println("cannot get users: ", err)
	}

	err = json.WriteJSON(w, http.StatusOK, &users)
	if err != nil {
		bookApp.Logger.ErrorLog.Println("Unable to write json to output: ", err)
	}
}
