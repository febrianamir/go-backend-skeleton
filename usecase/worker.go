package usecase

import (
	"app/lib/mailer"
	"app/lib/websocket"
	"app/request"
	"context"
)

func (usecase *Usecase) SendEmail(ctx context.Context, req request.SendEmailPayload) error {
	return usecase.repo.SendMailWithTemplate(ctx, mailer.SendMailWithTemplateParam{
		To:           req.To,
		TemplateName: req.TemplateName,
		TemplateData: req.TemplateData,
		Subject:      req.Subject,
	})
}

func (usecase *Usecase) BroadcastWebsocketMessage(ctx context.Context, message websocket.Message) error {
	return usecase.repo.BroadcastWebsocketMessage(ctx, message)
}
