package mail

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/ymanshur/simplebank/config"
)

func TestSendEmail_WithGmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	config, err := config.LoadConfig()
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "A test email"
	content := `
	<h1>Hello world</h1>
	<p>This is a test message from <a href="#">Simple Bank API</a></p>
	`
	to := []string{"ymanshur928@gmail.com"}
	attachFiles := []string{"../../README.md"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err)
}
