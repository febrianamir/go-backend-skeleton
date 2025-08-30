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
	log.Println(c)
}
