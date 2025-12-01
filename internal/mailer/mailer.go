package mailer

import (
	"bytes"
	"embed"
	"fmt"
	"text/template"

	"go.uber.org/zap"
)

const (
	templateDir 	  = "templates/"
	VerificationEmail = "verification.tmpl"
	ResetPassword     = "reset_password.tmpl"
	InvitationEmail   =  "invitation.tmpl"
)


var FS embed.FS


type EmailProvider interface {
	SendEmail(to, sbject, body string) error
}
type Client interface {
	Send(templateFile, email string, data any, isSandbox bool)(int error)
}


type EmailClient struct {
	Provider EmailProvider
	F5      embed.FS
	logger  *zap.SugaredLogger
}


func NewEmailClient(provider EmailProvider, logger *zap.SugaredLogger) *EmailClient{
	return &EmailClient{
		Provider: provider,
        F5: FS,
		logger: logger,
	}
}

func (e *EmailClient) Send(templateFile, email string, data any, isSandbox bool)(int, error){
	e.logger.Infof("Send email to: %s", email)
	tmpl, err := template.ParseFS(e.F5, templateDir+templateFile)
		if err != nil {
			return -1, err
		}
		
		subject := new(bytes.Buffer)
		err = tmpl.ExecuteTemplate(subject, "subjects", data)
		if err != nil {
			return -1, err
		}

		body := new(bytes.Buffer)
		err = tmpl.ExecuteTemplate(body, "body", data)
		if err != nil {
			return -1, err
		}


		if err := e.Provider.SendEmail(email, subject.String(), body.String()); err != nil {
			e.logger.Error("failed to send email", err)
			return 0, fmt.Errorf("failed to send email: %w", err)
}

e.logger.Infof("Email sent successfully to %s", email)

return 200, nil

}
