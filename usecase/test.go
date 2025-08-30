package usecase

import (
	"app/lib/mailer"
	"app/request"
	"context"
)

func (usecase *Usecase) TestSendEmail(ctx context.Context, req request.TestSendEmail) error {
	err := usecase.repo.SendMailWithTemplate(ctx, mailer.SendMailWithTemplateParam{
		To:           []string{req.Email},
		TemplateName: "test_email.html",
		TemplateData: map[string]any{
			"message": "Test Send Email",
		},
		Subject: "Test Send Email",
	})
	if err != nil {
		return err
	}

	return nil
}
