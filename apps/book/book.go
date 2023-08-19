package book

import (
	"go-base-structure/cmd/config"
)

var bookApp *config.Application

func NewBookApp(app *config.Application) {
	bookApp = app
}
