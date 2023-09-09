package book

import (
	"go-base-structure/cmd/settings"
)

// bookApp is wide configuration instance belong to book app
var bookApp *settings.Application

// NewBookApp assign sent wide configuration instance to the bookApp variable
func NewBookApp(app *settings.Application) {
	bookApp = app
}
