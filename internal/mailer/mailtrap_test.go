package mailer

import (
	"testing"

	"go.uber.org/zap"
)

const (
	templateFile = "verification.tmpl"
	username	 = "muslim_tech"
	useremail	 = "muslimtech@example.com"
	isSandbox	 = true
)


var data = map[string]string{
	"FirstName": "Yusuf",
	"Code": "123456",
}

var logger = zap.NewExample().Sugar()

func TestEmailClient_Send(t *testing.T){
	defer logger.Sync()

	mockProvider := &MockEmailProvider{ShouldFail: false}
	mailer := NewEmailClient(mockProvider, logger)

	status, err := mailer.Send(templateFile, useremail, data, isSandbox)
	if err != nil{
	t.Fatalf("Expected no error, go: %v", err)
	}
	if status != 200{
	t.Errorf("Expected status 200, go: %d", status)
	}
}


func TestEmailClient_Send_Failure(t *testing.T){
	mockProvider := &MockEmailProvider{ShouldFail: true}
	mailer := NewEmailClient(mockProvider, logger)	
	defer logger.Sync()

	_, err := mailer.Send(templateFile, useremail, data, isSandbox)
	if err == nil{
		t.Fatalf("Expected error, got nil")
	}

	if len(mockProvider.SentEmails) != 1 {
		t.Fatalf("Expected 1 email to be sent, got: %d", len(mockProvider.SentEmails))
	}

	email := mockProvider.SentEmails[0]
	if email.To != useremail {
		t.Errorf("Expected email to be sent to %s, got: %s", useremail, email.To)
	}
	if email.Subject == "" {	
		t.Error("Expected email subject to be non-empty")
	}
}