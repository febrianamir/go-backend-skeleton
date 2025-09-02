package usecase

import (
	"app/lib/mailer"
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
