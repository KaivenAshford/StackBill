package testutil

import (
	"testing"

	"github.com/kingqaquuu/stackbill/internal/middleware"
)

const TestJWTSecret = "test-secret-key"
const TestJWTExpire = 72

func GenerateTestToken(t *testing.T, userID uint, username string) string {
	t.Helper()
	token, err := middleware.GenerateToken(userID, username, TestJWTSecret, TestJWTExpire)
	if err != nil {
		t.Fatalf("generate test token: %v", err)
	}
	return token
}
