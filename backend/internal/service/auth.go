package service

import (
	"errors"
	"log/slog"

	"github.com/kingqaquuu/stackbill/internal/dto"
	"github.com/kingqaquuu/stackbill/internal/middleware"
	"github.com/kingqaquuu/stackbill/internal/model"
	"github.com/kingqaquuu/stackbill/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	userRepo     *repository.UserRepository
	categoryRepo *repository.CategoryRepository
	jwtSecret    string
	jwtExpire    int
}

func NewAuthService(userRepo *repository.UserRepository, categoryRepo *repository.CategoryRepository, jwtSecret string, jwtExpire int) *AuthService {
	return &AuthService{
		userRepo:     userRepo,
		categoryRepo: categoryRepo,
		jwtSecret:    jwtSecret,
		jwtExpire:    jwtExpire,
	}
}

func (s *AuthService) Register(req *dto.RegisterRequest) (*dto.LoginResponse, error) {
	if _, err := s.userRepo.FindByUsername(req.Username); err == nil {
		return nil, NewServiceError(409, ErrCodeDuplicateUsername, "username already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if _, err := s.userRepo.FindByEmail(req.Email); err == nil {
		return nil, NewServiceError(409, ErrCodeDuplicateEmail, "email already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	s.initDefaultCategories(user.ID)

	token, err := middleware.GenerateToken(user.ID, user.Username, s.jwtSecret, s.jwtExpire)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token: token,
		User:  dto.UserResponse{ID: user.ID, Username: user.Username, Email: user.Email, Nickname: user.Nickname, Avatar: user.Avatar},
	}, nil
}

func (s *AuthService) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, NewServiceError(401, ErrCodeInvalidCredentials, "invalid credentials")
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, NewServiceError(401, ErrCodeInvalidCredentials, "invalid credentials")
	}

	token, err := middleware.GenerateToken(user.ID, user.Username, s.jwtSecret, s.jwtExpire)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token: token,
		User:  dto.UserResponse{ID: user.ID, Username: user.Username, Email: user.Email, Nickname: user.Nickname, Avatar: user.Avatar},
	}, nil
}

func (s *AuthService) initDefaultCategories(userID uint) {
	categories := []model.Category{
		{UserID: userID, Name: "AI 工具", Type: "subscription", Color: "#7c3aed", Icon: "robot", SortOrder: 1},
		{UserID: userID, Name: "开发工具", Type: "subscription", Color: "#2563eb", Icon: "code", SortOrder: 2},
		{UserID: userID, Name: "云服务", Type: "subscription", Color: "#0891b2", Icon: "cloud", SortOrder: 3},
		{UserID: userID, Name: "域名", Type: "subscription", Color: "#059669", Icon: "globe", SortOrder: 4},
		{UserID: userID, Name: "服务器", Type: "subscription", Color: "#d97706", Icon: "server", SortOrder: 5},
		{UserID: userID, Name: "娱乐", Type: "subscription", Color: "#dc2626", Icon: "game-controller", SortOrder: 6},
		{UserID: userID, Name: "办公", Type: "subscription", Color: "#4f46e5", Icon: "briefcase", SortOrder: 7},
		{UserID: userID, Name: "其他", Type: "subscription", Color: "#6b7280", Icon: "ellipsis", SortOrder: 8},
	}
	if err := s.categoryRepo.BatchCreate(categories); err != nil {
		slog.Error("failed to init default categories", "user_id", userID, "error", err)
	}
}
