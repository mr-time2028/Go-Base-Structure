package main

import (
	"go-base-structure/cmd/server"
)

func main() {
	app, err := server.Serve()
	if err != nil {
		app.Logger.ErrorLog.Fatal(err)
	}
}
