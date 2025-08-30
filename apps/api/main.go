package main

import (
	"app"
	"app/config"
	"app/handler"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load()
}

func main() {
	config := config.InitConfig()
	db, err := config.NewDB()
	if err != nil {
		log.Fatal("failed connected to database: ", err)
	}

	app := app.NewApp(db)
	handler := handler.NewHandler(app)
	router := chi.NewRouter()

	router.Get("/healthz", handler.Healthz)
	router.Group(func(r chi.Router) {})

	serverAddr := fmt.Sprintf("0.0.0.0:%s", config.SERVER_PORT)
	server := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}
	log.Println("server listen at:", serverAddr)
	server.ListenAndServe()
}
