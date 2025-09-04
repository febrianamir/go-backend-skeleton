package usecase

import (
	"app/lib/constant"
	"app/request"
	"context"
)

func (usecase *Usecase) TestSendEmail(ctx context.Context, req request.TestSendEmail) error {
	return usecase.repo.PublishTask(ctx, constant.TaskTypeEmailSend, request.SendEmailPayload{
		To:           []string{req.Email},
		TemplateName: "test_email.html",
		TemplateData: map[string]any{
			"message": "Test Send Email",
		},
		Subject: "Test Send Email",
	})
}
