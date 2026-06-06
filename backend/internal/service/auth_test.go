package service

import (
	"testing"

	"github.com/kingqaquuu/stackbill/internal/dto"
)

func TestAuthService_Register_Success(t *testing.T) {
	svc := newTestServices(t)

	resp, err := svc.AuthService.Register(&dto.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	if resp.Token == "" {
		t.Error("expected token")
	}
	if resp.User.Username != "testuser" {
		t.Errorf("username = %q, want %q", resp.User.Username, "testuser")
	}
	if resp.User.Email != "test@example.com" {
		t.Errorf("email = %q, want %q", resp.User.Email, "test@example.com")
	}
}

func TestAuthService_Register_DuplicateUsername(t *testing.T) {
	svc := newTestServices(t)

	_, _ = svc.AuthService.Register(&dto.RegisterRequest{
		Username: "testuser", Email: "first@example.com", Password: "password123",
	})

	_, err := svc.AuthService.Register(&dto.RegisterRequest{
		Username: "testuser", Email: "second@example.com", Password: "password123",
	})
	if err == nil {
		t.Fatal("expected error for duplicate username")
	}
	svcErr, ok := err.(*ServiceError)
	if !ok {
		t.Fatalf("expected ServiceError, got %T", err)
	}
	if svcErr.Code != ErrCodeDuplicateUsername {
		t.Errorf("code = %d, want %d", svcErr.Code, ErrCodeDuplicateUsername)
	}
}

func TestAuthService_Register_DuplicateEmail(t *testing.T) {
	svc := newTestServices(t)

	_, _ = svc.AuthService.Register(&dto.RegisterRequest{
		Username: "user1", Email: "same@example.com", Password: "password123",
	})

	_, err := svc.AuthService.Register(&dto.RegisterRequest{
		Username: "user2", Email: "same@example.com", Password: "password123",
	})
	if err == nil {
		t.Fatal("expected error for duplicate email")
	}
	svcErr, ok := err.(*ServiceError)
	if !ok {
		t.Fatalf("expected ServiceError, got %T", err)
	}
	if svcErr.Code != ErrCodeDuplicateEmail {
		t.Errorf("code = %d, want %d", svcErr.Code, ErrCodeDuplicateEmail)
	}
}

func TestAuthService_Login_Success(t *testing.T) {
	svc := newTestServices(t)

	_, _ = svc.AuthService.Register(&dto.RegisterRequest{
		Username: "loginuser", Email: "login@example.com", Password: "password123",
	})

	resp, err := svc.AuthService.Login(&dto.LoginRequest{
		Username: "loginuser", Password: "password123",
	})
	if err != nil {
		t.Fatalf("login: %v", err)
	}
	if resp.Token == "" {
		t.Error("expected token after login")
	}
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	svc := newTestServices(t)

	_, _ = svc.AuthService.Register(&dto.RegisterRequest{
		Username: "loginuser2", Email: "login2@example.com", Password: "password123",
	})

	_, err := svc.AuthService.Login(&dto.LoginRequest{
		Username: "loginuser2", Password: "wrongpassword",
	})
	if err == nil {
		t.Fatal("expected error for wrong password")
	}
	svcErr, ok := err.(*ServiceError)
	if !ok {
		t.Fatalf("expected ServiceError, got %T", err)
	}
	if svcErr.Code != ErrCodeInvalidCredentials {
		t.Errorf("code = %d, want %d", svcErr.Code, ErrCodeInvalidCredentials)
	}
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	svc := newTestServices(t)

	_, err := svc.AuthService.Login(&dto.LoginRequest{
		Username: "nonexistent", Password: "password123",
	})
	if err == nil {
		t.Fatal("expected error for nonexistent user")
	}
}

func TestAuthService_Register_CreatesDefaultCategories(t *testing.T) {
	svc := newTestServices(t)

	resp, _ := svc.AuthService.Register(&dto.RegisterRequest{
		Username: "catuser", Email: "cat@example.com", Password: "password123",
	})

	cats, err := svc.CategoryService.List(resp.User.ID, &dto.CategoryListQuery{})
	if err != nil {
		t.Fatalf("list categories: %v", err)
	}
	if len(cats) != 8 {
		t.Errorf("got %d categories, want 8", len(cats))
	}
}
