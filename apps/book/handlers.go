package book

import (
	"go-base-structure/cmd/config"
	"go-base-structure/models"
	"go-base-structure/pkg/json"
	"net/http"
)

// Home is a simple handler for book app
func Home(w http.ResponseWriter, r *http.Request) {
	users, err := models.Models.Book.GetAll()
	if err != nil {
		config.AppConfig.ErrorLog.Println("cannot get users: ", err)
	}

	err = json.WriteJSON(w, http.StatusOK, &users)
	if err != nil {
		config.AppConfig.ErrorLog.Println("Unable to write json to output: ", err)
	}
}
