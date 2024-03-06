package utils

import (
	"fmt"
	"io"
	"log"
	"os"
)

func LogOnError(f func() error) {
	if err := f(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func ReadAsString(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("unable to open file: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatalf("unable to close file: %v", err)
		}
	}(file)

	// Read the SQL file
	content, err := io.ReadAll(file)
	return string(content), err
}
