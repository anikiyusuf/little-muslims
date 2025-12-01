package mailer



import (
	"errors"
)

type MockEmailProvider struct {
	ShouldFail		bool
	SentEmails		[]struct {
		To			string
		Subject		string
		Body        string
	}
}

func (m  *MockEmailProvider) SendEmail(to, subject, body string) error{
	if m.ShouldFail{
		return errors.New("simulated email send failure")
	}
	m.SentEmails = append(m.SentEmails, struct{
		To 		string
		Subject	string
		Body	string
    }{To: to, Subject:subject, Body:body})
	return nil
}