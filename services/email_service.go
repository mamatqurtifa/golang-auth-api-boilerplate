package services

import (
	"fmt"
	"net/smtp"

	"github.com/yourusername/golang-auth-api-boilerplate/config"
)

type EmailService struct{}

// NewEmailService creates a new email service
func NewEmailService() *EmailService {
	return &EmailService{}
}

// SendEmail sends an email
func (s *EmailService) SendEmail(to, subject, body string) error {
	cfg := config.AppConfig

	from := cfg.SMTPFrom
	password := cfg.SMTPPassword

	// SMTP server configuration
	smtpHost := cfg.SMTPHost
	smtpPort := fmt.Sprintf("%d", cfg.SMTPPort)

	// Message
	message := []byte(fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-version: 1.0;\r\n"+
			"Content-Type: text/html; charset=\"UTF-8\";\r\n"+
			"\r\n"+
			"%s\r\n",
		from, to, subject, body,
	))

	// Authentication
	auth := smtp.PlainAuth("", cfg.SMTPUsername, password, smtpHost)

	// Send email
	err := smtp.SendMail(
		smtpHost+":"+smtpPort,
		auth,
		from,
		[]string{to},
		message,
	)

	if err != nil {
		return err
	}

	return nil
}

// SendPasswordResetEmail sends a password reset email
func (s *EmailService) SendPasswordResetEmail(to, token string) error {
	cfg := config.AppConfig
	resetLink := fmt.Sprintf("%s/reset-password?token=%s", cfg.FrontendURL, token)

	subject := "Password Reset Request"
	body := fmt.Sprintf(`
		<html>
		<body>
			<h2>Password Reset Request</h2>
			<p>You have requested to reset your password. Click the link below to reset your password:</p>
			<p><a href="%s">Reset Password</a></p>
			<p>This link will expire in 1 hour.</p>
			<p>If you did not request this, please ignore this email.</p>
		</body>
		</html>
	`, resetLink)

	return s.SendEmail(to, subject, body)
}

// SendVerificationEmail sends an email verification email
func (s *EmailService) SendVerificationEmail(to, token string) error {
	cfg := config.AppConfig
	verificationLink := fmt.Sprintf("%s/verify-email?token=%s", cfg.FrontendURL, token)

	subject := "Email Verification"
	body := fmt.Sprintf(`
		<html>
		<body>
			<h2>Welcome!</h2>
			<p>Thank you for registering. Please verify your email address by clicking the link below:</p>
			<p><a href="%s">Verify Email</a></p>
			<p>If you did not register, please ignore this email.</p>
		</body>
		</html>
	`, verificationLink)

	return s.SendEmail(to, subject, body)
}
