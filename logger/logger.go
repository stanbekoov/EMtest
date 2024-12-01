package logger

import (
	"io"
	"log"
	"os"
)

var (
	INFO  = log.New(os.Stdout, "[INFO] ", log.Lshortfile|log.Ltime)
	DEBUG = log.New(os.Stdout, "[DEBUG] ", log.Lshortfile|log.Ltime)
	ERROR = log.New(os.Stdout, "[ERROR] ", log.Lshortfile|log.Ltime)
)

func SetOutputByFilename(logger *log.Logger, fileName string) *os.File {
	if fileName != "" {
		if fileName == "none" {
			logger.SetOutput(io.Discard)
			return nil
		}

		out, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Fatal("error opening infoStr. Check your .env inputs")
		}

		logger.SetOutput(out)
		return out
	}
	return nil
}
