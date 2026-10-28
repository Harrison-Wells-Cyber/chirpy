package auth

import "testing"

func TestHashAndPasswordCheck(t *testing.T) {
	t.Run("hash password", func(t *testing.T) {
		hash, err := HashPassword("test-password-123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if hash == "" {
			t.Fatal("expected a non-empty hash")
		}
		if hash == "test-password-123" {
			t.Fatal("hash should not equal the original password")
		}
	})

	t.Run("check matching password", func(t *testing.T) {
		hash, err := HashPassword("test-password-123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		match, err := CheckPasswordHash("test-password-123", hash)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !match {
			t.Fatal("expected matching password to return true")
		}
	})

	t.Run("check non-matching password", func(t *testing.T) {
		hash, err := HashPassword("test-password-123")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		match, err := CheckPasswordHash("wrong-password", hash)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if match {
			t.Fatal("expected non-matching password to return false")
		}
	})
}
