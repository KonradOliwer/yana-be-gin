package main

import (
	"open-note-ne-go/config"
)

func main() {
	app, db := config.SetupServer("resources/migrations")
	err := app.Run(":8000")
	if err != nil {
		panic(err)
	}
	db.Close()
}
