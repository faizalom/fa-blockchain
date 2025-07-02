package blockchain

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

var dbPath string

func InitDB() {
	dbPath = os.Getenv("BLOCK_DB")
	if dbPath == "" {
		log.Fatal("BLOCK_DB environment variable is not set")
	}
	fmt.Println(dbPath)
}

func GetData(filename string) ([]byte, error) {
	// Open the file for reading
	file, err := os.Open(filepath.Join(dbPath, filename))
	if err != nil {
		log.Println("Error opening file:", err)
		return nil, err
	}
	defer file.Close()

	// Read the entire file into memory
	data, err := io.ReadAll(file)
	if err != nil {
		log.Println("Error reading file:", err)
		return nil, err
	}
	return data, err
}

func SaveToFile(b []byte, filename string) {
	file, err := os.Create(filepath.Join(dbPath, filename))
	if err != nil {
		log.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	file.Write(b)
}

func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}
