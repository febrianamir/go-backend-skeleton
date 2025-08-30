package main

import (
	"app/config"
	"log"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func main() {
	c := config.InitConfig()
	db, err := c.NewDB()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(db)
}
