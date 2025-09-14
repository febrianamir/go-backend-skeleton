package config

import (
	"app/lib/websocket"
	"fmt"
)

func (c *Config) NewWebsocketPool(maxConns int) *websocket.WebsocketPool {
	wsUrl := fmt.Sprintf("%s?api_key=%s", c.WEBSOCKET_URL, c.WEBSOCKET_API_KEY)
	return websocket.NewWebsocketPool(wsUrl, maxConns)
}
