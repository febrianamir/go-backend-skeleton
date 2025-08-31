package config

import (
	"app/lib/storage"
	"log"
)

func NewLocalStorage() *storage.Local {
	log.Println("using local storage")
	return &storage.Local{Directory: "public"}
}
