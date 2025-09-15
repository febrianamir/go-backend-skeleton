package worker

import (
	"context"
	"encoding/json"

	"app/lib/websocket"
	"app/request"

	"github.com/hibiken/asynq"
)

func (w *Worker) WorkerSendEmail(ctx context.Context, t *asynq.Task) error {
	var p request.SendEmailPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}

	err := w.App.Usecase.SendEmail(ctx, p)
	if err != nil {
		return err
	}

	return nil
}

func (w *Worker) WorkerBroadcastWebsocketMessage(ctx context.Context, t *asynq.Task) error {
	var p websocket.Message
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}

	err := w.App.Usecase.BroadcastWebsocketMessage(ctx, p)
	if err != nil {
		return err
	}

	return nil
}
