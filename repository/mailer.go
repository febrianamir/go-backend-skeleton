package repository

import (
	"app/lib/logger"
	"app/lib/mailer"
	"app/lib/signoz"
	"context"

	"go.uber.org/zap"
)

func (repo *Repository) SendMailWithTemplate(ctx context.Context, param mailer.SendMailWithTemplateParam) error {
	ctx, span := signoz.StartSpan(ctx, "repository.SendMailWithTemplate")
	defer span.Finish()

	err := repo.mailer.SendMailWithTemplate(ctx, param)
	if err != nil {
		logger.LogError(ctx, "error send mail with template", []zap.Field{
			zap.Error(err),
			zap.Strings("tags", []string{"repository", "SendMailWithTemplate"}),
		}...)
		return err
	}
	return nil
}

func (repo *Repository) SendPlainMail(ctx context.Context, param mailer.SendPlainMailParam) error {
	ctx, span := signoz.StartSpan(ctx, "repository.SendPlainMail")
	defer span.Finish()

	err := repo.mailer.SendPlainMail(ctx, param)
	if err != nil {
		logger.LogError(ctx, "error send plain mail", []zap.Field{
			zap.Error(err),
			zap.Strings("tags", []string{"repository", "SendPlainMail"}),
		}...)
		return err
	}
	return nil
}
