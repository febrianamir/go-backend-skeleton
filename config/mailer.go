package config

import (
	"app/lib/mailer"

	"gopkg.in/gomail.v2"
)

func (c *Config) NewSMTP() *mailer.SMTP {
	client := gomail.NewDialer(c.SMTP_HOST, c.SMTP_PORT, c.SMTP_USERNAME, c.SMTP_PASSWORD)
	client.SSL = false

	return &mailer.SMTP{
		Client:         client,
		SMTPSender:     c.SMTP_SENDER,
		SMTPSenderName: c.SMTP_SENDER_NAME,
	}
}
