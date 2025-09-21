package repository

import (
	"app/lib/signoz"
	"app/lib/websocket"
	"context"
)

func (repo *Repository) BroadcastWebsocketMessage(ctx context.Context, message websocket.Message) error {
	ctx, span := signoz.StartSpan(ctx, "repository.BroadcastWebsocketMessage")
	defer span.Finish()

	return repo.wsPool.SendMessage(ctx, message)
}
