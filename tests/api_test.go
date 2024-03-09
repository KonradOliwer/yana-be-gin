package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"open-note-ne-go/config"
	"strconv"
	"strings"
	"testing"
)

func TestCreateNote(t *testing.T) {
	// Given
	app, db := config.SetupServer("../resources/migrations")
	defer dropDataFromNotesTable(db)

	//When
	requestBody := map[string]string{"name": "test name", "content": "test content"}
	responseRecorder := performRequest(app, "POST", "/notes/", &requestBody)

	//Then
	assert.Equal(t, 201, responseRecorder.Code)
	responseBody, err := bodyAsMap(responseRecorder)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, responseBody["name"], requestBody["name"])
	assert.Equal(t, responseBody["content"], requestBody["content"])
	_, ok := responseBody["id"]
	assert.True(t, ok)
}

func TestGetNotes(t *testing.T) {
	//Given
	app, db := config.SetupServer("../resources/migrations")
	defer dropDataFromNotesTable(db)

	requestBody := map[string]string{"name": "test name", "content": "test content"}
	responseRecorder := performRequest(app, "POST", "/notes/", &requestBody)

	assert.Equal(t, 201, responseRecorder.Code)
	_, err := bodyAsMap(responseRecorder)
	if err != nil {
		t.Fatal(err)
	}
	postResponseBody, err := bodyAsMap(responseRecorder)
	postResponseId := postResponseBody["id"]

	//When
	responseRecorder = performRequest(app, "GET", "/notes/", nil)

	//Then
	assert.Equal(t, 200, responseRecorder.Code)
	getResponseBody, err := bodyAsArrayOfMaps(responseRecorder)
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, len(getResponseBody) == 1)
	assert.Equal(t, getResponseBody[0]["id"], postResponseId)
	assert.Equal(t, getResponseBody[0]["name"], requestBody["name"])
	assert.Equal(t, getResponseBody[0]["content"], requestBody["content"])
}

func TestUpdateNote(t *testing.T) {
	//Given
	app, db := config.SetupServer("../resources/migrations")
	defer dropDataFromNotesTable(db)

	requestBody := map[string]string{"name": "test name", "content": "test content"}
	responseRecorder := performRequest(app, "POST", "/notes/", &requestBody)

	assert.Equal(t, 201, responseRecorder.Code)
	_, err := bodyAsMap(responseRecorder)
	if err != nil {
		t.Fatal(err)
	}
	postResponseBody, err := bodyAsMap(responseRecorder)
	if err != nil {
		t.Fatal(err)
	}
	postResponseId := postResponseBody["id"]

	//When
	requestBody = map[string]string{"id": fmt.Sprintf("%v", postResponseId), "name": "updated name", "content": "updated content"}
	responseRecorder = performRequest(app, "PUT", fmt.Sprintf("/notes/%v", postResponseId), &requestBody)

	//Then
	assert.Equal(t, 200, responseRecorder.Code)
	responseRecorder = performRequest(app, "GET", "/notes/", nil)

	assert.Equal(t, 200, responseRecorder.Code)
	getResponseBody, err := bodyAsArrayOfMaps(responseRecorder)
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, len(getResponseBody) == 1)
	assert.Equal(t, getResponseBody[0]["id"], postResponseId)
	assert.Equal(t, getResponseBody[0]["name"], requestBody["name"])
	assert.Equal(t, getResponseBody[0]["content"], requestBody["content"])
}

func TestDeleteNotes(t *testing.T) {
	//Given
	app, db := config.SetupServer("../resources/migrations")
	defer dropDataFromNotesTable(db)

	requestBody := map[string]string{"name": "test name", "content": "test content"}
	responseRecorder := performRequest(app, "POST", "/notes/", &requestBody)

	assert.Equal(t, 201, responseRecorder.Code)
	postResponseBody, err := bodyAsMap(responseRecorder)
	if err != nil {
		t.Fatal(err)
	}
	postResponseId := postResponseBody["id"]

	//When
	responseRecorder = performRequest(app, "DELETE", fmt.Sprintf("/notes/%v", postResponseId), nil)

	//Then
	assert.Equal(t, 204, responseRecorder.Code)

	responseRecorder = performRequest(app, "GET", "/notes/", nil)
	assert.Equal(t, 200, responseRecorder.Code)
	getResponseBody, err := bodyAsArrayOfMaps(responseRecorder)
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, len(getResponseBody) == 0)
}

func TestGetMultipleNotes(t *testing.T) {
	//Given
	app, db := config.SetupServer("../resources/migrations")
	defer dropDataFromNotesTable(db)

	for i := 0; i < 3; i++ {
		requestBody := map[string]string{"name": "test name" + strconv.Itoa(i), "content": "test content" + strconv.Itoa(i)}
		responseRecorder := performRequest(app, "POST", "/notes/", &requestBody)

		assert.Equal(t, 201, responseRecorder.Code)
		_, err := bodyAsMap(responseRecorder)
		if err != nil {
			t.Fatal(err)
		}
		_, err = bodyAsMap(responseRecorder)
		if err != nil {
			t.Fatal(err)
		}
	}

	//When
	responseRecorder := performRequest(app, "GET", "/notes/", nil)

	//Then
	assert.Equal(t, 200, responseRecorder.Code)
	getResponseBody, err := bodyAsArrayOfMaps(responseRecorder)
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, len(getResponseBody) == 3)
}

func TestCreationOfNotesWithSameName(t *testing.T) {
	// Given
	app, db := config.SetupServer("../resources/migrations")
	defer dropDataFromNotesTable(db)

	requestBody := map[string]string{"name": "test name", "content": "test content"}
	responseRecorder := performRequest(app, "POST", "/notes/", &requestBody)
	assert.Equal(t, 201, responseRecorder.Code)

	//When
	requestBody = map[string]string{"name": "test name", "content": "test content2"}
	responseRecorder = performRequest(app, "POST", "/notes/", &requestBody)

	//Then
	assert.Equal(t, 400, responseRecorder.Code)
	responseBody, err := bodyAsMap(responseRecorder)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, responseBody["code"], "NOTE_ALREADY_EXISTS")
}

func TestCreationNoteWithTooLongName(t *testing.T) {
	// Given
	app, db := config.SetupServer("../resources/migrations")
	defer dropDataFromNotesTable(db)

	requestBody := map[string]string{"name": strings.Repeat("a", 51), "content": "test content"}
	responseRecorder := performRequest(app, "POST", "/notes/", &requestBody)
	assert.Equal(t, 400, responseRecorder.Code)
	postResponseBody, err := bodyAsMap(responseRecorder)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, postResponseBody["code"], "NOTE_VALIDATION_ERROR")
	assert.Equal(t, postResponseBody["message"], "name: the length of name must be between 1 and 50")
}

func TestUpdateWithTooLongName(t *testing.T) { //Given
	app, db := config.SetupServer("../resources/migrations")
	defer dropDataFromNotesTable(db)

	requestBody := map[string]string{"name": "test name", "content": "test content"}
	responseRecorder := performRequest(app, "POST", "/notes/", &requestBody)

	assert.Equal(t, 201, responseRecorder.Code)
	_, err := bodyAsMap(responseRecorder)
	if err != nil {
		t.Fatal(err)
	}
	postResponseBody, err := bodyAsMap(responseRecorder)
	if err != nil {
		t.Fatal(err)
	}
	noteId := fmt.Sprintf("%v", postResponseBody["id"])

	//When
	requestBody = map[string]string{"id": noteId, "name": strings.Repeat("a", 51), "content": "updated content"}
	responseRecorder = performRequest(app, "PUT", "/notes/"+noteId, &requestBody)

	//Then
	assert.Equal(t, 400, responseRecorder.Code)
	postResponseBody, err = bodyAsMap(responseRecorder)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, postResponseBody["code"], "NOTE_VALIDATION_ERROR")
	assert.Equal(t, postResponseBody["message"], "name: the length of name must be between 1 and 50")
}

func TestAttemptUpdateOnNonExistingEntity(t *testing.T) {
	// Given
	app, db := config.SetupServer("../resources/migrations")
	defer dropDataFromNotesTable(db)

	id := "b7460307-6f64-4e4a-b57f-2114478b48ba" // not existing valid id
	requestBody := map[string]string{"id": id, "name": "test name", "content": "test content"}
	responseRecorder := performRequest(app, "PUT", "/notes/"+id, &requestBody)
	assert.Equal(t, 404, responseRecorder.Code)
	postResponseBody, err := bodyAsMap(responseRecorder)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, postResponseBody["code"], "NOTE_NOT_FOUND")
}

func TestUpdateWithIdNotMatchingUrl(t *testing.T) { //Given
	app, db := config.SetupServer("../resources/migrations")
	defer dropDataFromNotesTable(db)

	requestBody := map[string]string{"name": "test name", "content": "test content"}
	responseRecorder := performRequest(app, "POST", "/notes/", &requestBody)

	assert.Equal(t, 201, responseRecorder.Code)
	_, err := bodyAsMap(responseRecorder)
	if err != nil {
		t.Fatal(err)
	}
	postResponseBody, err := bodyAsMap(responseRecorder)
	if err != nil {
		t.Fatal(err)
	}
	noteId := fmt.Sprintf("%v", postResponseBody["id"])

	//When
	requestBody = map[string]string{"id": "7b5a4832-de96-448d-bd38-b2b127c8f409", "name": "updated name", "content": "updated content"}
	responseRecorder = performRequest(app, "PUT", "/notes/"+noteId, &requestBody)

	//Then
	assert.Equal(t, 400, responseRecorder.Code)
	postResponseBody, err = bodyAsMap(responseRecorder)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, postResponseBody["code"], "NOTE_VALIDATION_ERROR")
	assert.Equal(t, postResponseBody["message"], "id: should match url id")
}

func bodyAsMap(w *httptest.ResponseRecorder) (map[string]interface{}, error) {
	var bodyData map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &bodyData)
	return bodyData, err
}

func bodyAsArrayOfMaps(w *httptest.ResponseRecorder) ([]map[string]interface{}, error) {
	var bodyData []map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &bodyData)
	return bodyData, err
}

func dropDataFromNotesTable(db *sql.DB) {
	_, err := db.Exec("DELETE FROM notes")
	if err != nil {
		return
	}
}

func performRequest(app *gin.Engine, method string, url string, body *map[string]string) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	if body != nil {
		err, serialisedBody := createRequestBody(*body)
		if err != nil {
			panic(err)
		}
		request, err := http.NewRequest(method, url, serialisedBody)
		if err != nil {
			panic(err)
		}
		app.ServeHTTP(recorder, request)
		return recorder
	} else {
		request, err := http.NewRequest(method, url, nil)
		if err != nil {
			panic(err)
		}
		app.ServeHTTP(recorder, request)
		return recorder
	}
}

func createRequestBody(body map[string]string) (error, *bytes.Buffer) {
	data := body
	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	serialisedBody := bytes.NewBuffer(jsonData)
	return err, serialisedBody
}
