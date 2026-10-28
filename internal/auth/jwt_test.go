package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWT(t *testing.T) {
	secret := "test-token-secret"
	userID := uuid.New()

	t.Run("make and validate jwt", func(t *testing.T) {
		tokenString, err := MakeJWT(userID, secret, 10*time.Minute)
		if err != nil {
			t.Fatalf("unexpected error making jwt: %v", err)
		}
		if tokenString == "" {
			t.Fatal("expected a non-empty token string")
		}

		gotID, err := ValidateJWT(tokenString, secret)
		if err != nil {
			t.Fatalf("unexpected error validating jwt: %v", err)
		}
		if gotID != userID {
			t.Fatalf("expected user id %s, got %s", userID, gotID)
		}
	})

	t.Run("expired token rejected", func(t *testing.T) {
		tokenString, err := MakeJWT(userID, secret, -1*time.Second)
		if err != nil {
			t.Fatalf("unexpected error making jwt: %v", err)
		}
		if _, err := ValidateJWT(tokenString, secret); err == nil {
			t.Fatal("expected error for expired token, got nil")
		}
	})

	t.Run("wrong secret rejected", func(t *testing.T) {
		tokenString, err := MakeJWT(userID, secret, 10*time.Minute)
		if err != nil {
			t.Fatalf("unexpected error making jwt: %v", err)
		}
		if _, err := ValidateJWT(tokenString, "different-secret"); err == nil {
			t.Fatal("expected error for wrong secret, got nil")
		}
	})
}
