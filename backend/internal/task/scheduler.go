package task

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/kingqaquuu/stackbill/internal/config"
	"github.com/kingqaquuu/stackbill/internal/repository"
	"github.com/kingqaquuu/stackbill/pkg/email"
)

type Scheduler struct {
	notificationRepo *repository.NotificationRepository
	subscriptionRepo *repository.SubscriptionRepository
	assetRepo        *repository.AssetRepository
	userRepo         *repository.UserRepository
	emailSender      *email.Sender
	stopCh           chan struct{}
}

func NewScheduler(
	notificationRepo *repository.NotificationRepository,
	subscriptionRepo *repository.SubscriptionRepository,
	assetRepo *repository.AssetRepository,
	userRepo *repository.UserRepository,
	smtpCfg config.SMTPConfig,
) *Scheduler {
	return &Scheduler{
		notificationRepo: notificationRepo,
		subscriptionRepo: subscriptionRepo,
		assetRepo:        assetRepo,
		userRepo:         userRepo,
		emailSender:      email.NewSender(smtpCfg),
		stopCh:           make(chan struct{}),
	}
}

func (s *Scheduler) Start() {
	if !s.emailSender.Enabled() {
		slog.Info("SMTP not configured, reminder email scheduler disabled")
		return
	}

	slog.Info("starting reminder email scheduler")
	go s.run()
}

func (s *Scheduler) Stop() {
	close(s.stopCh)
}

func (s *Scheduler) run() {
	// Run once on start, then daily at 9am
	s.checkAndSend()
	for {
		now := time.Now()
		next := time.Date(now.Year(), now.Month(), now.Day()+1, 9, 0, 0, 0, now.Location())
		timer := time.NewTimer(next.Sub(now))
		select {
		case <-timer.C:
			s.checkAndSend()
		case <-s.stopCh:
			timer.Stop()
			return
		}
	}
}

func (s *Scheduler) checkAndSend() {
	slog.Info("running daily reminder email check")

	settings, err := s.notificationRepo.GetAllEmailEnabled()
	if err != nil {
		slog.Error("failed to get notification settings", "error", err)
		return
	}

	for _, setting := range settings {
		user, err := s.userRepo.FindByID(setting.UserID)
		if err != nil {
			continue
		}

		var reminders []string

		// Check subscriptions renewing soon
		days := setting.RemindDaysBefore
		if days <= 0 {
			days = 3
		}
		subs, _ := s.subscriptionRepo.GetActiveByUserID(setting.UserID)
		deadline := time.Now().AddDate(0, 0, days)
		for _, sub := range subs {
			if sub.NextPaymentDate != nil && !sub.NextPaymentDate.After(deadline) {
				reminders = append(reminders, fmt.Sprintf("- %s: renewal on %s (%.2f %s)",
					sub.Name, sub.NextPaymentDate.Format("2006-01-02"), sub.Amount, sub.Currency))
			}
		}

		// Check assets expiring soon
		assets, _ := s.assetRepo.GetByUserID(setting.UserID)
		for _, asset := range assets {
			if asset.ExpireDate != nil && !asset.ExpireDate.After(deadline) {
				reminders = append(reminders, fmt.Sprintf("- %s (%s): expires on %s",
					asset.Name, asset.AssetType, asset.ExpireDate.Format("2006-01-02")))
			}
		}

		if len(reminders) == 0 {
			continue
		}

		body := fmt.Sprintf("You have %d upcoming reminder(s):\n\n%s\n\n-- StackBill",
			len(reminders), fmt.Sprintf("%s", reminders))
		subject := fmt.Sprintf("StackBill: %d upcoming reminder(s)", len(reminders))

		if err := s.emailSender.Send(user.Email, subject, body); err != nil {
			slog.Error("failed to send reminder email", "user_id", setting.UserID, "error", err)
		} else {
			slog.Info("sent reminder email", "user_id", setting.UserID, "count", len(reminders))
		}
	}
}
