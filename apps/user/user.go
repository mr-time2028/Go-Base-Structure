package user

import (
	"go-base-structure/cmd/settings"
)

// userApp is wide configuration instance belong to user app
var userApp *settings.Application

// NewUserApp assign sent wide configuration instance to the userApp variable
func NewUserApp(app *settings.Application) {
	userApp = app
}
