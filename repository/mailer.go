package repository

import (
	"app/lib/mailer"
	"context"
)

func (repo *Repository) SendMailWithTemplate(ctx context.Context, param mailer.SendMailWithTemplateParam) error {
	return repo.mailer.SendMailWithTemplate(ctx, param)
}

func (repo *Repository) SendPlainMail(ctx context.Context, param mailer.SendPlainMailParam) error {
	return repo.mailer.SendPlainMail(ctx, param)
}
