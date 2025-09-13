package websocket

import "time"

// Message represents a message
type Message struct {
	MessageType  string        `json:"message_type"`
	Notification *Notification `json:"notification"`
	Timestamp    time.Time     `json:"timestamp"`
}

// Notification represents a notification message
type Notification struct {
	NotificationType string `json:"notification_type"`
	Title            string `json:"title"`
	Message          string `json:"message"`
	Level            string `json:"level"`
}
