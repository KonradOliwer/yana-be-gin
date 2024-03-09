package notes

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/google/uuid"
	"open-note-ne-go/database"
)

type Service struct {
	Db *sql.DB
}

type NoteView struct {
	Id      string `json:"id"`
	Name    string `json:"name" binding:"required,max=50"`
	Content string `json:"content"`
}

type ErrorView struct {
	Code    ErrorResponseCodes `json:"code"`
	Message string             `json:"message"`
}

type ErrorResponseCodes string

const (
	NOT_FOUND        ErrorResponseCodes = "NOTE_NOT_FOUND"
	ALREADY_EXISTS   ErrorResponseCodes = "NOTE_ALREADY_EXISTS"
	VALIDATION_ERROR ErrorResponseCodes = "NOTE_VALIDATION_ERROR"
	UNKNOWN_ERROR    ErrorResponseCodes = "NOTE_UNKNOWN_ERROR"
)

func NewNoteView(note database.Note) NoteView {
	return NoteView{
		Id:      note.Id,
		Name:    note.Name,
		Content: note.Content,
	}
}

func (service *Service) RegisterRoutes(router *gin.Engine) {
	notesGroup := router.Group("/notes")
	notesGroup.Use(ErrorHandler)
	notesGroup.GET("/", service.getAll)
	notesGroup.POST("/", service.add)
	notesGroup.PUT("/:noteId", service.update)
	notesGroup.DELETE("/:noteId", service.delete)
}

func ErrorHandler(c *gin.Context) {
	c.Next()

	// TODO: this only handles the first error, we should handle all of them
	for _, err := range c.Errors {
		switch err.Error() {
		case "pq: duplicate key value violates unique constraint \"notes_name_key\"":
			c.JSON(400, ErrorView{Code: ALREADY_EXISTS})
			return
		//This should be more generic in future
		case "Key: 'NoteView.Name' Error:Field validation for 'Name' failed on the 'max' tag":
			c.JSON(400, ErrorView{Code: VALIDATION_ERROR, Message: "name: the length of name must be between 1 and 50"})
			return
		case "id: should match url id":
			c.JSON(400, ErrorView{Code: VALIDATION_ERROR, Message: "id: should match url id"})
			return
		case "sql: no rows in result set":
			c.JSON(404, ErrorView{Code: NOT_FOUND})
			return
		default:
			c.JSON(500, ErrorView{Code: UNKNOWN_ERROR})
			return
		}
	}
}

func (service *Service) getAll(c *gin.Context) {
	notes, err := database.GetAllNotes(service.Db)
	if err != nil {
		c.Error(err)
		return
	}
	noteViews := make([]NoteView, len(notes))
	for i, note := range notes {
		noteViews[i] = NewNoteView(note)
	}
	c.JSON(200, noteViews)
}

func (service *Service) add(c *gin.Context) {
	var inputNote NoteView
	if err := c.ShouldBindJSON(&inputNote); err != nil {
		c.Error(err)
		return
	}

	note, err := database.CreateNote(service.Db, inputNote.Name, inputNote.Content)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(201, NewNoteView(note))
}

func (service *Service) update(c *gin.Context) {
	var inputNote NoteView
	if err := c.ShouldBindJSON(&inputNote); err != nil {
		c.Error(err)
		return
	}
	if inputNote.Id != c.Param("noteId") {
		c.Error(fmt.Errorf("id: should match url id"))
		return
	}

	id, err := uuid.Parse(inputNote.Id)
	note, err := database.UpdateNote(service.Db, id, inputNote.Name, inputNote.Content)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, NewNoteView(note))
}

func (service *Service) delete(c *gin.Context) {
	noteId := c.Param("noteId")
	parse, err := uuid.Parse(noteId)
	if err != nil {
		c.Error(err)
		return
	}
	err = database.DeleteById(service.Db, parse)
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(204)
}
