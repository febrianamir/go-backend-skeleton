package main

import (
	"app"
	"app/config"
	"app/handler"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
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

	app := app.NewApp(cfg, db, mailer, storage, cache, publisher)
	handler := handler.NewHandler(app)
	router := chi.NewRouter()

	router.Get("/healthz", handler.Healthz)
	router.Group(func(r chi.Router) {
		r.Use(handler.InstrumentMiddleware)

		// Test
		r.Route("/tests", func(r chi.Router) {
			r.Post("/send-email", handler.TestSendEmail)
		})

		// File
		r.Route("/files", func(r chi.Router) {
			r.Post("/upload", handler.UploadFile)
		})

		// Auth
		r.Route("/auth", func(r chi.Router) {
			r.Route("/register", func(r chi.Router) {
				r.Post("/", handler.Register)
				r.Post("/resend-verification", handler.RegisterResendVerification)
				r.Post("/verify-account", handler.VerifyAccount)
			})
			r.Post("/login", handler.Login)
		})

		// User
		r.Route("/users", func(r chi.Router) {
			r.Get("/", handler.GetUsers)
			r.Post("/", handler.CreateUser)
			r.Get("/{ID}", handler.GetUser)
			r.Put("/{ID}", handler.UpdateUser)
			r.Delete("/{ID}", handler.DeleteUser)
		})
	})

	serverAddr := fmt.Sprintf("0.0.0.0:%s", cfg.SERVER_PORT)
	server := &http.Server{
		Addr:         serverAddr,
		WriteTimeout: time.Duration(cfg.SERVER_WRITE_TIMEOUT) * time.Second,
		ReadTimeout:  time.Duration(cfg.SERVER_READ_TIMEOUT) * time.Second,
		IdleTimeout:  time.Duration(cfg.SERVER_IDLE_TIMEOUT) * time.Second,
		Handler:      router,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt)

	go func() {
		log.Println("server listen at:", serverAddr)
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed to start: %v", err)
		}
	}()

	<-quit
	log.Println("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.SERVER_SHUTDOWN_TIMEOUT)*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("server forced to shutdown: %v", err)
		if closeErr := server.Close(); closeErr != nil {
			log.Printf("error closing server: %v", closeErr)
		}
	}

	log.Println("server gracefully stopped")
}
