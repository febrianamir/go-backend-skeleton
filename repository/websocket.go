package repository

import (
	"app/lib/logger"
	"app/lib/websocket"
	"context"
	"fmt"

	coderWs "github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
	"go.uber.org/zap"
)

func (repo *Repository) BroadcastWebsocketMessage(ctx context.Context, message websocket.Message) error {
	url := fmt.Sprintf("%s?api_key=%s", repo.config.WEBSOCKET_URL, repo.config.WEBSOCKET_API_KEY)
	c, _, err := coderWs.Dial(ctx, url, nil)
	if err != nil {
		logger.LogError(ctx, "error websocket.Dial", []zap.Field{
			zap.Error(err),
			zap.Strings("tags", []string{"repository", "BroadcastWebsocketMessage"}),
		}...)
		return err
	}
	defer c.CloseNow()

	err = wsjson.Write(ctx, c, message)
	if err != nil {
		logger.LogError(ctx, "error wsjson.Write", []zap.Field{
			zap.Error(err),
			zap.Strings("tags", []string{"repository", "BroadcastWebsocketMessage"}),
		}...)
		return err
	}

	err = c.Close(coderWs.StatusNormalClosure, "")
	if err != nil {
		logger.LogError(ctx, "error conn.Close", []zap.Field{
			zap.Error(err),
			zap.Strings("tags", []string{"repository", "BroadcastWebsocketMessage"}),
		}...)
		return err
	}
	return nil
}
