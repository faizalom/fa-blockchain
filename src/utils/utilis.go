package utils

import (
	"bytes"
	"encoding/gob"
	"flag"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func SetEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file:", err)
	}

	blockDb := flag.String("block_db", os.Getenv("BLOCK_DB"), "Blockchain database directory path")
	errLogFile := flag.String("error_log_file", os.Getenv("ERROR_LOG_FILE"), "Error log file path")
	deviceId := flag.String("device_id", os.Getenv("DEVICE_ID"), "device id used for msq broker connector id")
	sqlDb := flag.String("sql_db", os.Getenv("SQL_DB"), "SQL database connection string")

	flag.Parse()

	if *deviceId == "" {
		log.Println("Device ID is not set, using a new UUID")
		*deviceId = uuid.NewString()
	}
	os.Setenv("BLOCK_DB", *blockDb)
	os.Setenv("ERROR_LOG_FILE", *errLogFile)
	os.Setenv("DEVICE_ID", *deviceId)
	os.Setenv("SQL_DB", *sqlDb)
}

func NewEventId() string {
	eventIdv7, err := uuid.NewV7()
	if err != nil {
		log.Println(err)
	}
	return eventIdv7.String()
}

func Serialize(data any) ([]byte, error) {
	var buffer bytes.Buffer

	encode := gob.NewEncoder(&buffer)
	err := encode.Encode(data)
	return buffer.Bytes(), err
}
