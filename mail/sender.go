package mail

import (
	"context"
	"fmt"
	
	"github.com/wneessen/go-mail"
)

const (
	smtpGmailHost = "smtp.gmail.com"
	smtpGmailPort = 587
)

type EmailSender interface {
	SendEmail(
		subject string,
		body string,
		to []string,
		cc []string,
		bcc []string,
		attachFiles []string,
	) error
	Ping(ctx context.Context) error
}

type GmailSender struct {
	name             string
	fromEmailAddress string
	client           *mail.Client
}

func NewGmailSender(name, fromEmailAddress, fromEmailPassword string) (EmailSender, error) {
	client, err := mail.NewClient(smtpGmailHost, mail.WithPort(smtpGmailPort),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(fromEmailAddress), mail.WithPassword(fromEmailPassword))
	if err != nil {
		return nil, err
	}
	
	return &GmailSender{
		name:             name,
		fromEmailAddress: fromEmailAddress,
		client:           client,
	}, nil
}

func (sender *GmailSender) SendEmail(
	subject string,
	body string,
	to []string,
	cc []string,
	bcc []string,
	attachFiles []string,
) error {
	message := mail.NewMsg()
	
	if err := message.FromFormat(sender.name, sender.fromEmailAddress); err != nil {
		return fmt.Errorf("failed to set our From address: %w", err)
	}
	
	message.Subject(subject)
	message.SetBodyString(mail.TypeTextHTML, body)
	
	if err := message.To(to...); err != nil {
		return fmt.Errorf("failed to set To addresses: %w", err)
	}
	
	if err := message.Cc(cc...); err != nil {
		return fmt.Errorf("failed to set Cc addresses: %w", err)
	}
	
	if err := message.Bcc(bcc...); err != nil {
		return fmt.Errorf("failed to set Bcc addresses: %w", err)
	}
	
	for _, file := range attachFiles {
		message.AttachFile(file)
	}
	
	if err := sender.client.DialAndSend(message); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	
	return nil
}

// Ping checks if the email client is connected to the SMTP server.
func (sender *GmailSender) Ping(ctx context.Context) error {
	return sender.client.DialWithContext(ctx)
}
