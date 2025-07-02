package utils

import (
	"log"
	"os"
)

func LogErrors(errorLogFile string) {
	if errorLogFile != "" {
		logFile, err := os.OpenFile(errorLogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening file: %v", err)
		}

		log.SetOutput(logFile)
		log.SetFlags(log.LstdFlags | log.Llongfile)
		// } else {
		// 	log.SetFlags(0)
	}
	log.SetFlags(log.LstdFlags | log.Llongfile)
}
