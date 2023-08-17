package book

import (
	"go-base-structure/cmd/config"
	"go-base-structure/helpers"
	"go-base-structure/models"
	"net/http"
)

// Home is a simple handler for book app
func Home(w http.ResponseWriter, r *http.Request) {
	users, err := models.Models.Book.GetAll()
	if err != nil {
		config.AppConfig.ErrorLog.Println("cannot get users: ", err)
	}

	err = helpers.WriteJSON(w, http.StatusOK, &users)
	if err != nil {
		config.AppConfig.ErrorLog.Println("Unable to write json to output: ", err)
	}
}
