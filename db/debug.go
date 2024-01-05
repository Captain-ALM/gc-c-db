package db

import (
	"log"
	"os"
)

func DebugPrintln(msg string) {
	if os.Getenv("DEBUG_DB") == "1" {
		log.Println("DEBUG_DB:", msg)
	}
}

func DebugErrIsNil(err error) bool {
	if err == nil {
		return true
	}
	DebugPrintln(err.Error())
	return false
}
