package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"open-note-ne-go/utils"
	"os"
)

type ConnectionData struct {
	Host     string
	Port     int
	DbName   string
	User     string
	Password string
	SslMode  string
}

func DefaultPGAdminConnection() ConnectionData {
	return ConnectionData{
		Host:     "localhost",
		Port:     5432,
		DbName:   "yana",
		User:     "user",
		Password: "password",
		SslMode:  "disable",
	}
}

const operationDbName = "yana"

func InitialiseDatabase(migrationFolder string) (*sql.DB, error) {
	db, err := initialiseDatabaseConnection(DefaultPGAdminConnection())
	if err != nil {
		return nil, err
	}

	//TODO create db if not exists and then connect to new db and close old connection

	err = migrate(db, migrationFolder)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func initialiseDatabaseConnection(connection ConnectionData) (*sql.DB, error) {
	var passwordString string
	if connection.Password == "" {
		passwordString = ""
	} else {
		passwordString = " password=" + connection.Password
	}
	connectionParameters := fmt.Sprintf("host=%s port=%d user=%s%s dbname=%s sslmode=%s",
		connection.Host, connection.Port, connection.User, passwordString, connection.DbName, connection.SslMode)
	db, err := sql.Open("postgres", connectionParameters)
	if err != nil {
		return nil, fmt.Errorf("error while opening the database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error while pinging the database: %v", err)
	}
	fmt.Println("Successfully connected!")
	return db, nil
}

func migrate(db *sql.DB, migrationFolder string) error {
	//Here should be robust migration system run, but for now we will only run all the filed in migrations folder sorted by name
	files, err := os.ReadDir(migrationFolder)
	if err != nil {
		return fmt.Errorf("error while reading migrations folder: %v", err)
	}
	for _, file := range files {
		query, err := utils.ReadAsString(migrationFolder + "/" + file.Name())
		_, err = db.Exec(query)
		if err != nil {
			return fmt.Errorf("error while executing migration file %s: %v", file.Name(), err)
		}
	}
	return nil
}
