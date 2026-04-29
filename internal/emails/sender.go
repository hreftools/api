package emails

import "context"

type EmailSendParams struct {
	To      []string
	Html    string
	Text    string
	Subject string
}

type EmailSender interface {
	Send(ctx context.Context, params EmailSendParams) error
}
