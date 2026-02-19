package emails

import (
	"github.com/resend/resend-go/v3"
)

type EmailSender interface {
	Send(params *resend.SendEmailRequest) error
}

type ResendEmailSender struct {
	client *resend.Client
}

func (s *ResendEmailSender) Send(params *resend.SendEmailRequest) error {
	_, err := s.client.Emails.Send(params)
	return err
}

func NewResendEmailSender(client *resend.Client) *ResendEmailSender {
	return &ResendEmailSender{
		client: client,
	}
}
