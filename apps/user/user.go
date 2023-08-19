package user

import (
	"go-base-structure/cmd/config"
)

var userApp *config.Application

func NewUserApp(app *config.Application) {
	userApp = app
}
