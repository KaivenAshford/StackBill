package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/kingqaquuu/stackbill/internal/dto"
	"github.com/kingqaquuu/stackbill/internal/model"
	"github.com/kingqaquuu/stackbill/internal/repository"
	"gorm.io/gorm"
)

type WebhookService struct {
	repo *repository.WebhookRepository
}

func NewWebhookService(repo *repository.WebhookRepository) *WebhookService {
	return &WebhookService{repo: repo}
}

func (s *WebhookService) List(userID uint) ([]dto.WebhookResponse, error) {
	webhooks, err := s.repo.List(userID)
	if err != nil {
		return nil, err
	}
	result := make([]dto.WebhookResponse, len(webhooks))
	for i, w := range webhooks {
		result[i] = dto.WebhookResponse{ID: w.ID, URL: w.URL, Events: w.Events, Active: w.Active}
	}
	return result, nil
}

func (s *WebhookService) Create(userID uint, req *dto.CreateWebhookRequest) (*dto.WebhookResponse, error) {
	webhook := &model.Webhook{
		UserID: userID,
		URL:    req.URL,
		Secret: req.Secret,
		Events: req.Events,
		Active: true,
	}
	if err := s.repo.Create(webhook); err != nil {
		return nil, err
	}
	return &dto.WebhookResponse{ID: webhook.ID, URL: webhook.URL, Events: webhook.Events, Active: webhook.Active}, nil
}

func (s *WebhookService) Update(userID uint, id uint, req *dto.UpdateWebhookRequest) (*dto.WebhookResponse, error) {
	webhook, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewServiceError(404, ErrCodeNotFound, "webhook not found")
		}
		return nil, err
	}
	if webhook.UserID != userID {
		return nil, NewServiceError(403, ErrCodeForbidden, "forbidden")
	}
	if req.URL != "" {
		webhook.URL = req.URL
	}
	webhook.Secret = req.Secret
	if req.Events != "" {
		webhook.Events = req.Events
	}
	if req.Active != nil {
		webhook.Active = *req.Active
	}
	if err := s.repo.Update(webhook); err != nil {
		return nil, err
	}
	return &dto.WebhookResponse{ID: webhook.ID, URL: webhook.URL, Events: webhook.Events, Active: webhook.Active}, nil
}

func (s *WebhookService) Delete(userID uint, id uint) error {
	webhook, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return NewServiceError(404, ErrCodeNotFound, "webhook not found")
		}
		return err
	}
	if webhook.UserID != userID {
		return NewServiceError(403, ErrCodeForbidden, "forbidden")
	}
	return s.repo.Delete(id)
}

// Deliver sends a webhook event to all active webhooks for the user.
func (s *WebhookService) Deliver(userID uint, event string, payload map[string]any) {
	webhooks, err := s.repo.GetActiveByUserID(userID)
	if err != nil {
		slog.Error("failed to get webhooks for delivery", "user_id", userID, "error", err)
		return
	}

	body, _ := json.Marshal(map[string]any{
		"event":     event,
		"payload":   payload,
		"timestamp": time.Now().Format(time.RFC3339),
	})

	for _, w := range webhooks {
		if !strings.Contains(w.Events, event) {
			continue
		}
		go s.sendWebhook(w, body)
	}
}

func (s *WebhookService) sendWebhook(w model.Webhook, body []byte) {
	req, err := http.NewRequest("POST", w.URL, strings.NewReader(string(body)))
	if err != nil {
		slog.Error("webhook request creation failed", "url", w.URL, "error", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	if w.Secret != "" {
		mac := hmac.New(sha256.New, []byte(w.Secret))
		mac.Write(body)
		sig := hex.EncodeToString(mac.Sum(nil))
		req.Header.Set("X-Webhook-Signature", fmt.Sprintf("sha256=%s", sig))
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("webhook delivery failed", "url", w.URL, "error", err)
		return
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)

	if resp.StatusCode >= 300 {
		slog.Error("webhook returned error status", "url", w.URL, "status", resp.StatusCode)
	}
}
