package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"app"
	"app/config"
	"app/scheduler"

	"github.com/go-co-op/gocron/v2"
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

	app := app.NewApp(cfg, db, mailer, storage, cache, publisher, nil)

	// create a scheduler
	s, err := scheduler.NewScheduler(app)
	if err != nil {
		log.Fatal("failed create new scheduler: ", err)
	}

	// add a job to the scheduler
	s.RegisterJob(gocron.DurationJob(60*time.Second), "CronTest", s.App.Usecase.CronTest)

	// start the scheduler
	s.Start()

	// scheduler gracefull shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt)

	<-quit
	log.Println("shutting down scheduler...")

	// when you're done, shut it down
	err = s.Shutdown()
	if err != nil {
		log.Fatalf("failed to shutdown scheduler, do forced shutdown: %v", err)
	}

	log.Println("scheduler gracefully stopped")
}
