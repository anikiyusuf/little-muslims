package mailer 

import (
	"fmt"
	"log"
	"net/smtp"

	
	"github.com/jordan-wright/email"
)



type MailtrapProvider struct {
	Host	string
	Port	string
	User	string
	Pass	string 
    From    string
}

func NewMailtrapProvider(host, user, pass, from, port string) *MailtrapProvider{
	return &MailtrapProvider{
		Host:	host,
		Port:	port,
		User:	user,
		Pass:	pass,
		From:	from,
	}
}

func (m *MailtrapProvider) SendEmail(to, subject, body string) error{
	email := email.NewEmail()
	email.From	= fmt.Sprintf("Sender name <%s>", m.From)
	email.To	= []string{to}
	email.Subject = subject
	email.HTML = []byte(body)

		// SMTP configuration
		smtpAuth := smtp.PlainAuth("", m.User, m.Pass, m.Host)
		err :=  email.Send(fmt.Sprintf("%s:%s", m.Host, m.Port), smtpAuth)
		if err != nil {
			log.Fatalf("Failed to send email:%v", err)
		}

		return nil
}