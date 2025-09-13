package websocket

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"app"
	"app/lib/auth"

	"github.com/coder/websocket"
)

type Websocket struct {
	App *app.App
}

func NewWebsocket(a *app.App) *Websocket {
	return &Websocket{App: a}
}

// WebSocket handler
func (ws *Websocket) HandleWebSocket(h *Hub, w http.ResponseWriter, r *http.Request) {
	idTokenClaims := auth.GetAuthFromCtx(r.Context())
	if idTokenClaims.Subject == "" {
		log.Printf("Failed to authenticate user")
		return
	}

	// Accept WebSocket connection
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true, // For development only, set to false in production
		// OriginPatterns:     []string{"*"}, // Allow all origins for testing
	})
	if err != nil {
		log.Printf("Failed to accept WebSocket connection: %v", err)
		return
	}

	// Create context for this connection, don't use context from request
	ctx, cancel := context.WithCancel(context.Background())

	// Create client
	client := &Client{
		conn:   conn,
		send:   make(chan Message, 256),
		hub:    h,
		id:     fmt.Sprintf("client-%d-%s", time.Now().UnixNano(), idTokenClaims.Subject),
		userId: idTokenClaims.UserID,
		ctx:    ctx,
		cancel: cancel,
	}

	// Register client
	client.hub.register <- client

	// Start goroutines, to send and receive messages from client
	go client.writePump()
	go client.readPump()
}
