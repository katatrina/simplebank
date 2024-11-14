package mail

import (
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
}

type GmailSender struct {
	name              string
	fromEmailAddress  string
	fromEmailPassword string
}

func NewGmailSender(name, fromEmailAddress, fromEmailPassword string) EmailSender {
	return &GmailSender{
		name:              name,
		fromEmailAddress:  fromEmailAddress,
		fromEmailPassword: fromEmailPassword,
	}
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
	
	client, err := mail.NewClient(smtpGmailHost, mail.WithPort(smtpGmailPort),
		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithUsername(sender.fromEmailAddress), mail.WithPassword(sender.fromEmailPassword))
	if err != nil {
		return fmt.Errorf("failed to establish our email client: %w", err)
	}
	
	if err = client.DialAndSend(message); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	
	return nil
}
