package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"open-note-ne-go/database"
)

func main() {
	db, err := database.InitialiseDatabase("resources/migrations")
	if err != nil {
		panic(err)
	}

	ns, err := database.GetAllNotes(db)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Number of rows returned by the query: %d\n", len(ns))
	for _, note := range ns {
		fmt.Printf("Id: %s, Name: %s, Content: ```%s```\n", note.Id, note.Name, note.Content)
	}
}
