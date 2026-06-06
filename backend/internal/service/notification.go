package service

import (
	"errors"

	"github.com/kingqaquuu/stackbill/internal/dto"
	"github.com/kingqaquuu/stackbill/internal/model"
	"github.com/kingqaquuu/stackbill/internal/repository"
	"gorm.io/gorm"
)

type NotificationService struct {
	repo *repository.NotificationRepository
}

func NewNotificationService(repo *repository.NotificationRepository) *NotificationService {
	return &NotificationService{repo: repo}
}

func (s *NotificationService) Get(userID uint) (*dto.NotificationSettingResponse, error) {
	setting, err := s.repo.GetByUserID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &dto.NotificationSettingResponse{EmailEnabled: false, RemindDaysBefore: 3}, nil
		}
		return nil, err
	}
	return &dto.NotificationSettingResponse{
		ID:               setting.ID,
		EmailEnabled:     setting.EmailEnabled,
		RemindDaysBefore: setting.RemindDaysBefore,
	}, nil
}

func (s *NotificationService) Update(userID uint, req *dto.UpdateNotificationSettingRequest) (*dto.NotificationSettingResponse, error) {
	setting, err := s.repo.GetByUserID(userID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		setting = &model.NotificationSetting{
			UserID:          userID,
			EmailEnabled:    false,
			RemindDaysBefore: 3,
		}
		if err := s.repo.Create(setting); err != nil {
			return nil, err
		}
	}

	if req.EmailEnabled != nil {
		setting.EmailEnabled = *req.EmailEnabled
	}
	if req.RemindDaysBefore != nil {
		setting.RemindDaysBefore = *req.RemindDaysBefore
	}

	if err := s.repo.Update(setting); err != nil {
		return nil, err
	}

	return &dto.NotificationSettingResponse{
		ID:               setting.ID,
		EmailEnabled:     setting.EmailEnabled,
		RemindDaysBefore: setting.RemindDaysBefore,
	}, nil
}
