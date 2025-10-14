package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	"github.com/ymanshur/simplebank/config"
	"github.com/ymanshur/simplebank/internal/repo"
	"github.com/ymanshur/simplebank/internal/server/worker/presentation"
	"github.com/ymanshur/simplebank/internal/server/worker/tasktype"
	"github.com/ymanshur/simplebank/internal/ucase"
	"github.com/ymanshur/simplebank/pkg/mail"
)

type taskSendVerifyEmail struct {
	config config.Config
	repo   repo.Repo
	mailer mail.EmailSender

	verifyEmail ucase.VerifyEmailUcase
}

func NewTaskSendVerifyEmail(
	config config.Config,
	repo repo.Repo,
	mailer mail.EmailSender,
	verifyEmail ucase.VerifyEmailUcase,
) Task {
	return &taskSendVerifyEmail{
		config:      config,
		repo:        repo,
		mailer:      mailer,
		verifyEmail: verifyEmail,
	}
}

func (t *taskSendVerifyEmail) Name() string {
	return tasktype.SendVerifyEmail
}

func (t *taskSendVerifyEmail) Handler(ctx context.Context, payload []byte) error {
	var p presentation.SendVerifyEmailPayload
	if err := json.Unmarshal(payload, &p); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	verifyEmail, err := t.verifyEmail.Create(ctx, ucase.CreateVerifyEmailRequest{
		Username: p.Username,
	})
	if err != nil {
		return fmt.Errorf("failed to create verify email: %w", err)
	}

	subject := "Welcome to Simple Bank"
	// TODO: replace gRPC Gateway URL with an environment variable that points to a front-end page
	serverAddress := strings.Split(t.config.GRPCGatewayServerAddress, ":")
	serverPort := serverAddress[1]
	verifyUrl := fmt.Sprintf("http://localhost:%s/api/verify_user?email_id=%d&secret_code=%s",
		serverPort, verifyEmail.ID, verifyEmail.SecretCode)
	content := fmt.Sprintf(`Hello %s,<br/>
	Thank you for registering with us!<br/>
	Please <a href="%s">click here</a> to verify your email address.<br/>
	`, verifyEmail.FullName, verifyUrl)
	to := []string{verifyEmail.Email}

	err = t.mailer.SendEmail(subject, content, to, nil, nil, nil)
	if err != nil {
		return fmt.Errorf("failed to send verify email: %w", err)
	}

	log.Info().
		Str("type", t.Name()).
		Bytes("payload", payload).
		Str("email", verifyEmail.Email).
		Msg("processed task")
	return nil
}
