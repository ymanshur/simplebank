package mail

import (
	"github.com/ymanshur/simplebank/pkg/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSendEmail_WithGmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	config, err := util.LoadConfig("../..")
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "A test email"
	content := `
	<h1>Hello world</h1>
	<p>This is a test message from <a href="#">Simple Bank API</a></p>
	`
	to := []string{"yusuf.manshur@privy.id"}
	var attachFiles []string

	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err)
}
