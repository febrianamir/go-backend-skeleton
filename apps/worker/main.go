package main

import (
	"app"
	"app/config"
	"app/lib/constant"
	"app/worker"
	"log"

	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"
)

var cfg *config.Config

func init() {
	godotenv.Load()
	cfg = config.InitConfig()
	cfg.NewLogger()
}

func main() {
	db, err := cfg.NewDB()
	if err != nil {
		log.Fatal("failed connect to database: ", err)
	}
	mailer := cfg.NewSMTP()
	storage := cfg.NewStorage()
	redis := cfg.NewRedis()
	publisher, err := cfg.NewPublisher()
	if err != nil {
		log.Fatal("failed connect to publisher: ", err)
	}

	app := app.NewApp(db, mailer, storage, redis, publisher)
	worker := worker.NewWorker(app)
	server := cfg.NewConsumer()

	mux := asynq.NewServeMux()
	mux.HandleFunc(constant.TaskTypeEmailSend, worker.WorkerSendEmail)

	if err := server.Run(mux); err != nil {
		log.Fatalf("consumer server failed to start: %v", err)
	}
}
