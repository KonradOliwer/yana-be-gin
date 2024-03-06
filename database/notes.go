package database

import (
	"database/sql"
	"fmt"
	"open-note-ne-go/utils"
)

type Note struct {
	Id      string
	Name    string
	Content string
}

func GetAllNotes(db *sql.DB) ([]Note, error) {
	query := "SELECT id, name, content FROM notes"
	rows, err := db.Query(query)
	if err != nil {
		return []Note{}, fmt.Errorf("error while executing the query [(%s)]:\n%v", query, err)
	}
	defer utils.LogOnError(func() error { return rows.Close() })

	var notes []Note
	for i := 0; rows.Next(); i++ {
		var note Note
		if err := rows.Scan(&note.Id, &note.Name, &note.Content); err != nil {
			return []Note{}, fmt.Errorf("error while mapping the row[%d] from query [(%s)]:\n%v", i, query, err)
		}
		notes = append(notes, note)
	}
	return notes, nil
}

func CreateNote(db *sql.DB, note Note) (*Note, error) {
	_, err := db.Exec("INSERT INTO notes (id, name, content) VALUES ($1, $2, $3)", note.Id, note.Name, note.Content)
	if err != nil {
		return &Note{}, fmt.Errorf("error while creatring note with : %v", err)
	}
	return &note, nil
}

func UpdateByName(db *sql.DB, note Note) (*Note, error) {
	_, err := db.Exec("UPDATE notes SET content=$1 WHERE name=$2", note.Content, note.Name)
	if err != nil {
		return &Note{}, fmt.Errorf("error while updating existring note: %v", err)
	}
	return &note, nil
}

func DeleteByName(db *sql.DB, name string) error {
	_, err := db.Exec("DELETE FROM notes WHERE name=$1", name)
	if err != nil {
		return fmt.Errorf("error while deleting note: %v", err)
	}
	return nil
}
