package emails

type EmailSendParams struct {
	From    string
	To      []string
	Html    string
	Text    string
	Subject string
	ReplyTo string
}

type EmailSender interface {
	Send(params EmailSendParams) error
}
