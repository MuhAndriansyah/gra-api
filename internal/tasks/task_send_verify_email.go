package tasks

import (
	"backend-layout/internal/adapter/mail"
	"backend-layout/internal/config"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	"gopkg.in/gomail.v2"
)

const (
	TaskSendVerifyEmail = "task:send_verify_email"
)

type PayloadSendVerifyEmail struct {
	Email      string
	Username   string
	VerifyCode string
}

func (r *RedisTaskDestributor) DistributeTaskSendVerifyEmail(ctx context.Context, payload *PayloadSendVerifyEmail, opts ...asynq.Option) error {
	// serialisasi
	jsonPayload, err := json.Marshal(payload)

	if err != nil {
		return fmt.Errorf("failed to marshal task payload %w", err)
	}

	task := asynq.NewTask(TaskSendVerifyEmail, jsonPayload)
	taskInfo, err := r.client.EnqueueContext(ctx, task, opts...)

	if err != nil {
		log.Error().
			Err(err).
			Str("email", payload.Email).
			Str("username", payload.Username).
			Msg("failed to enqueue task")
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Info().
		Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Str("queue", taskInfo.Queue).
		Int("max_retry", taskInfo.MaxRetry).
		Msg("enqueued task")

	return nil
}

func HandlerVerifyEmail(ctx context.Context, task *asynq.Task) error {
	var payload PayloadSendVerifyEmail
	// deserialisasi
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	d := mail.InitMail()
	m := gomail.NewMessage()
	emailVerificationBody := BuldTemplateVerifyEmail(payload.Username, payload.VerifyCode)

	m.SetHeader("From", config.LoadMailConfig().MailEmail)
	m.SetHeader("To", payload.Email)
	m.SetHeader("Subject", "Email Verification")
	m.SetBody("text/html", emailVerificationBody)

	if err := d.DialAndSend(m); err != nil {
		log.Error().Err(err).Msg("failed to send verify email")
		return fmt.Errorf("failed to send email to %s: %w", payload.Email, err)
	}

	log.Info().Msg("email delivery task completed successfully")
	return nil
}

func BuldTemplateVerifyEmail(username, verifyCode string) string {
	absolutePath, _ := os.Getwd()

	filename := filepath.Join(absolutePath, "/internal/adapter/mail/template/", "email-verify.templ")

	tmpl, err := template.ParseFiles(filename)

	if err != nil {
		log.Error().Err(err).Msg("failed to parse file email verification templ")
	}

	payload := struct {
		Username   string
		VerifyCode string
	}{
		Username:   username,
		VerifyCode: verifyCode,
	}

	var out bytes.Buffer

	err = tmpl.Execute(&out, payload)

	if err != nil {
		log.Error().Err(err).Msg("failed to execute template and payload email verification")
	}

	return out.String()
}
