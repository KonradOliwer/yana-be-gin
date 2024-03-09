package database

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"open-note-ne-go/utils"
)

type Note struct {
	Id      string
	Name    string
	Content string
}

type Scanner interface {
	Scan(dest ...any) error
}

func GetAllNotes(db *sql.DB) ([]Note, error) {
	query := "SELECT id, name, content FROM notes"
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error while executing the query [(%s)]:\n%v", query, err)
	}
	defer utils.LogOnError(func() error { return rows.Close() })

	// In case of no results we want function to return empty slice not nil, do we need  `= []Note{}`
	var notes []Note = []Note{}
	for i := 0; rows.Next(); i++ {
		note, err2 := toNote(rows)
		if err2 != nil {
			return nil, err2
		}
		notes = append(notes, note)
	}
	return notes, nil
}

func CreateNote(db *sql.DB, name string, content string) (Note, error) {
	query := "INSERT INTO notes (id, name, content) VALUES ($1, $2, $3) RETURNING *"
	result := db.QueryRow(query, uuid.New(), name, content)
	note, err := toNote(result)
	if err != nil {
		return Note{}, err
	}
	return note, nil
}

func UpdateNote(db *sql.DB, id uuid.UUID, name string, content string) (Note, error) {
	query := "UPDATE notes SET name=$2, content=$3 WHERE id=$1 RETURNING *"
	result := db.QueryRow(query, id, name, content)
	note, err := toNote(result)
	if err != nil {
		return Note{}, err
	}
	return note, nil
}

func DeleteById(db *sql.DB, id uuid.UUID) error {
	_, err := db.Exec("DELETE FROM notes WHERE id=$1", id)
	if err != nil {
		return fmt.Errorf("error while deleting note: %v", err)
	}
	return nil
}

func toNote(row Scanner) (Note, error) {
	var note Note
	if err := row.Scan(&note.Id, &note.Name, &note.Content); err != nil {
		return Note{}, err
	}
	return note, nil
}
