package database

import (
	"testing"
)

func TestCreatedNoteIsAvailableViaGetAll(t *testing.T) {
	db, err := InitialiseDatabase("../resources/migrations")
	if err != nil {
		t.Errorf("Error while initialising the database: %v", err)
	}
	_, err = CreateNote(db, Note{
		Id:      "4548e374-6f94-4496-929d-64577f68b71a",
		Name:    "Test note",
		Content: "This is a test note",
	})
	if err != nil {
		t.Errorf("Error while creating the note: %v", err)
	}
	defer DeleteByName(db, "Test note")

	allNotes, err := GetAllNotes(db)
	if err != nil {
		t.Errorf("Error while getting all notes: %v", err)
	}

	if len(allNotes) != 1 {
		t.Errorf("Expected 1 note, got %d", len(allNotes))
	}

	if allNotes[0].Id != "4548e374-6f94-4496-929d-64577f68b71a" {
		t.Errorf("Expected note with id 4548e374-6f94-4496-929d-64577f68b71a, got %s", allNotes[0].Id)
	}

	if allNotes[0].Name != "Test note" {
		t.Errorf("Expected note with name Test note, got %s", allNotes[0].Name)
	}

	if allNotes[0].Content != "This is a test note" {
		t.Errorf("Expected note with content This is a test note, got %s", allNotes[0].Content)
	}
}

func TestUpdatedNoteIsAvailableViaGetALL(t *testing.T) {
	db, err := InitialiseDatabase("../resources/migrations")
	if err != nil {
		t.Errorf("Error while initialising the database: %v", err)
	}
	_, err = CreateNote(db, Note{
		Id:      "4548e374-6f94-4496-929d-64577f68b71a",
		Name:    "Test note",
		Content: "This is a test note",
	})
	if err != nil {
		t.Errorf("Error while creating the note: %v", err)
	}

	allNotes, err := GetAllNotes(db)
	if err != nil {
		t.Errorf("Error while getting all notes: %v", err)
	}

	if len(allNotes) != 1 {
		t.Errorf("Expected 1 note, got %d", len(allNotes))
	}

	if allNotes[0].Id != "4548e374-6f94-4496-929d-64577f68b71a" {
		t.Errorf("Expected note with id 4548e374-6f94-4496-929d-64577f68b71a, got %s", allNotes[0].Id)
	}

	if allNotes[0].Name != "Test note" {
		t.Errorf("Expected note with name Test note, got %s", allNotes[0].Name)
	}

	if allNotes[0].Content != "This is a test note" {
		t.Errorf("Expected note with content This is a test note, got %s", allNotes[0].Content)
	}
}
