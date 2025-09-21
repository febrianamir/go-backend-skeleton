package usecase

import (
	"app/lib/constant"
	"app/lib/signoz"
	"app/lib/websocket"
	"app/request"
	"context"
	"time"
)

func (usecase *Usecase) TestSendNotification(ctx context.Context, req request.TestSendNotification) error {
	ctx, span := signoz.StartSpan(ctx, "usecase.TestSendNotification")
	defer span.Finish()

	return usecase.repo.PublishTask(ctx, constant.TaskTypeWebsocketBroadcastMessage, websocket.Message{
		MessageType: websocket.MessageTypeNotification,
		Notification: &websocket.Notification{
			NotificationType: "TEST",
			Title:            req.Title,
			Message:          req.Message,
		},
		Timestamp: time.Now(),
	})
}
