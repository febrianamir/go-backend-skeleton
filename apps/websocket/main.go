package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"app"
	"app/config"
	"app/handler"
	"app/lib/websocket"

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

	app := app.NewApp(cfg, db, mailer, storage, cache, publisher, nil)
	handler := handler.NewHandler(app)

	// Create and start websocket hub
	ws := websocket.NewWebsocket()
	go ws.Hub.Run()

	// Set up HTTP routes
	router := chi.NewRouter()

	// WebSocket routes with authentication
	router.Route("/ws", func(r chi.Router) {
		r.With(handler.WebSocketAuthMiddleware).Get("/", func(w http.ResponseWriter, r *http.Request) {
			ws.HandleWebSocket(w, r)
		})
	})

	// Start server
	port := cfg.SERVER_WEBSOCKET_PORT
	serverAddr := fmt.Sprintf("0.0.0.0:%s", port)
	server := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt)

	go func() {
		log.Printf("websocket notification server starting on port %s", port)
		log.Printf("websocket endpoint: ws://localhost%s/ws", port)
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("websocket server failed to start: %v", err)
		}
	}()

	<-quit
	log.Println("shutting down websocket server...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.WEBSOCKET_SERVER_SHUTDOWN_TIMEOUT)*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("websocket server forced to shutdown: %v", err)
		if closeErr := server.Close(); closeErr != nil {
			log.Printf("error closing websocket server: %v", closeErr)
		}
	}

	log.Println("websocket server gracefully stopped")
}
