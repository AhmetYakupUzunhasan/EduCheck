package middleware

import (
	"fmt"
	"net/smtp"
	"strconv"
	"strings"
)

type GmailMailer struct {
	fromEmail string
	appPass   string
}

func NewGmailMailer(fromEmail, appPassword string) *GmailMailer {
	return &GmailMailer{
		fromEmail: fromEmail,
		appPass:   appPassword,
	}
}

// Send sends a simple UTF-8 text email via Gmail SMTP (STARTTLS on 587).
//
// Requirements (Gmail):
// - Use an App Password (recommended). Normal account password usually won't work.
// - The Gmail account must have 2FA enabled to create an App Password.
func (m *GmailMailer) Send(toEmail string, subject string, body int) error {
	toEmail = strings.TrimSpace(toEmail)
	if toEmail == "" {
		return fmt.Errorf("toEmail is required")
	}
	if strings.TrimSpace(subject) == "" {
		return fmt.Errorf("subject is required")
	}

	const host = "smtp.gmail.com"
	const addr = "smtp.gmail.com:587"

	auth := smtp.PlainAuth("", m.fromEmail, m.appPass, host)

	bodyStr := strconv.Itoa(body)

	msg := buildTextMessage(m.fromEmail, toEmail, subject, bodyStr)

	if err := smtp.SendMail(addr, auth, m.fromEmail, []string{toEmail}, []byte(msg)); err != nil {
		return fmt.Errorf("send mail: %w", err)
	}
	return nil
}

func buildTextMessage(from, to, subject, body string) string {
	// Minimal headers + blank line + body. Use CRLF per RFC.
	var b strings.Builder
	b.WriteString("From: " + from + "\r\n")
	b.WriteString("To: " + to + "\r\n")
	b.WriteString("Subject: " + sanitizeHeader(subject) + "\r\n")
	b.WriteString("MIME-Version: 1.0\r\n")
	b.WriteString("Content-Type: text/plain; charset=utf-8\r\n")
	b.WriteString("Content-Transfer-Encoding: 8bit\r\n")
	b.WriteString("\r\n")
	b.WriteString(body)
	if !strings.HasSuffix(body, "\n") {
		b.WriteString("\r\n")
	}
	return b.String()
}

func sanitizeHeader(s string) string {
	// Prevent header injection by stripping CR/LF.
	s = strings.ReplaceAll(s, "\r", "")
	s = strings.ReplaceAll(s, "\n", "")
	return s
}
