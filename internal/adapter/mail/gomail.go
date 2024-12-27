package mail

import (
	"backend-layout/internal/config"

	"gopkg.in/gomail.v2"
)

func InitMail() *gomail.Dialer {

	dialer := gomail.NewDialer(
		config.LoadMailConfig().SMTPHost,
		config.LoadMailConfig().SMTPPort,
		config.LoadMailConfig().MailUsername,
		config.LoadMailConfig().MailPassword,
	)

	return dialer
}
