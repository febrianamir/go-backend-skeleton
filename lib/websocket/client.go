package websocket

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/coder/websocket"
)

// Client represents a connected WebSocket client
type Client struct {
	conn      *websocket.Conn
	send      chan Message
	hub       *Hub
	id        string
	userId    uint
	userRoles []string
	ctx       context.Context
	cancel    context.CancelFunc
}

// writePump pumps messages from the hub to the websocket connection
func (c *Client) writePump() {
	defer func() {
		c.conn.Close(websocket.StatusInternalError, "write pump closed")
		c.cancel()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.Close(websocket.StatusNormalClosure, "channel closed")
				return
			}

			if err := c.conn.Write(c.ctx, websocket.MessageText, mustMarshal(message)); err != nil {
				log.Printf("Error writing to client %s: %v", c.id, err)
				return
			}

		case <-c.ctx.Done():
			return
		}
	}
}

// readPump pumps messages from the websocket connection to the hub, and handle client disconnect
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.cancel()
	}()

	for {
		_, message, err := c.conn.Read(c.ctx)
		if err != nil {
			if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
				log.Printf("Client %s disconnected normally", c.id)
			} else {
				log.Printf("Error reading from client %s: %v", c.id, err)
			}
			break
		}

		log.Printf("Received message from client %s: %s", c.id, string(message))

		// Parse and broadcast message
		clientIdSplits := strings.Split(c.id, "-")
		if len(clientIdSplits) == 3 && clientIdSplits[2] == "server" {
			// Parse and broadcast message from server to server connection
			var msg Message
			if err := json.Unmarshal(message, &msg); err == nil {
				select {
				case c.hub.broadcast <- msg:
				default:
					log.Printf("Failed to broadcast message to client %s", c.id)
				}
			} else {
				log.Printf("Failed to parse message: %s\n", err.Error())
			}
		}

		// Write parse and broadcast message from frontend to server connection here...
	}
}
