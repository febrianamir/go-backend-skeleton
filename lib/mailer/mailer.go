package mailer

import (
	"bytes"
	"context"
	"fmt"
	"html/template"

	"gopkg.in/gomail.v2"
)

type SMTP struct {
	Client         *gomail.Dialer
	SMTPSender     string
	SMTPSenderName string
}

func (s *SMTP) SendMailWithTemplate(ctx context.Context, param SendMailWithTemplateParam) error {
	mailMessage := gomail.NewMessage()
	mailMessage.SetHeader("From", s.SMTPSender)
	mailMessage.SetHeader("To", param.To...)
	if len(param.Cc) > 0 {
		mailMessage.SetHeader("Cc", param.Cc...)
	}
	mailMessage.SetHeader("Subject", param.Subject)

	t, err := template.ParseFiles(fmt.Sprintf("./asset/email/%s", param.TemplateName))
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	err = t.Execute(buf, param.TemplateData)
	if err != nil {
		return err
	}
	mailMessage.SetBody("text/html", buf.String())

	err = s.Client.DialAndSend(mailMessage)
	if err != nil {
		return err
	}

	return nil
}

func (s *SMTP) SendPlainMail(ctx context.Context, param SendPlainMailParam) error {
	mailMessage := gomail.NewMessage()
	mailMessage.SetHeader("From", s.SMTPSender)
	mailMessage.SetHeader("To", param.To...)
	if len(param.Cc) > 0 {
		mailMessage.SetHeader("Cc", param.Cc...)
	}
	mailMessage.SetHeader("Subject", param.Subject)
	mailMessage.SetBody("text/plain", param.Message)

	err := s.Client.DialAndSend(mailMessage)
	if err != nil {
		return err
	}

	return nil
}

type SendMailWithTemplateParam struct {
	To           []string       `json:"to"`
	Cc           []string       `json:"cc"`
	TemplateName string         `json:"template_name"`
	TemplateData map[string]any `json:"template_data"`
	Subject      string         `json:"subject"`
}

type SendPlainMailParam struct {
	To      []string `json:"to"`
	Cc      []string `json:"cc"`
	Subject string   `json:"subject"`
	Message string   `json:"message"`
}
