package email

import (
	"fmt"
	"log/slog"
	"net/smtp"
	"strings"

	"github.com/kingqaquuu/stackbill/internal/config"
)

type Sender struct {
	cfg config.SMTPConfig
}

func NewSender(cfg config.SMTPConfig) *Sender {
	return &Sender{cfg: cfg}
}

func (s *Sender) Enabled() bool {
	return s.cfg.Host != "" && s.cfg.User != ""
}

func (s *Sender) Send(to, subject, body string) error {
	if !s.Enabled() {
		slog.Warn("email not configured, skipping send")
		return nil
	}

	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
	from := s.cfg.From
	if from == "" {
		from = s.cfg.User
	}

	msg := strings.Join([]string{
		"From: " + from,
		"To: " + to,
		"Subject: " + subject,
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		"",
		body,
	}, "\r\n")

	auth := smtp.PlainAuth("", s.cfg.User, s.cfg.Password, s.cfg.Host)
	return smtp.SendMail(addr, auth, from, []string{to}, []byte(msg))
}
