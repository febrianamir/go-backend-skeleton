package repository

import (
	"app/lib/websocket"
	"context"
)

func (repo *Repository) BroadcastWebsocketMessage(ctx context.Context, message websocket.Message) error {
	return repo.wsPool.SendMessage(ctx, message)
}
