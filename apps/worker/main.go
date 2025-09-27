package main

import (
	"log"

	"app"
	"app/config"
	"app/lib/constant"
	"app/worker"

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
	cache, err := cfg.NewCache()
	if err != nil {
		log.Fatal("failed connect to redis: ", err)
	}
	publisher, err := cfg.NewPublisher()
	if err != nil {
		log.Fatal("failed connect to publisher: ", err)
	}
	wsPool := cfg.NewWebsocketPool(20)

	app := app.NewApp(cfg, db, mailer, storage, cache, publisher, wsPool)
	worker := worker.NewWorker(app)
	server := cfg.NewConsumer()

	mux := asynq.NewServeMux()
	worker.RegisterWorker(mux, constant.TaskTypeEmailSend, "WorkerSendEmail", false, worker.WorkerSendEmail)
	worker.RegisterWorker(mux, constant.TaskTypeWebsocketBroadcastMessage, "WorkerBroadcastWebsocketMessage", false, worker.WorkerBroadcastWebsocketMessage)

	if err := server.Run(mux); err != nil {
		log.Fatalf("consumer server failed to start: %v", err)
	}
}
