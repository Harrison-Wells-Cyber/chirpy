package auth

import (
	"encoding/hex"
	"net/http"
	"testing"
)

func TestMakeRefreshToken(t *testing.T) {
	t.Run("returns a 64-char hex string", func(t *testing.T) {
		token := MakeRefreshToken()
		if len(token) != 64 {
			t.Fatalf("expected 64 hex chars (32 bytes), got %d", len(token))
		}
		if _, err := hex.DecodeString(token); err != nil {
			t.Fatalf("expected a valid hex string, got %q: %v", token, err)
		}
	})

	t.Run("tokens are unique", func(t *testing.T) {
		a := MakeRefreshToken()
		b := MakeRefreshToken()
		if a == b {
			t.Fatal("expected two different tokens, got the same one twice")
		}
	})
}

func TestGetBearerToken(t *testing.T) {
	t.Run("valid bearer token", func(t *testing.T) {
		headers := http.Header{}
		headers.Set("Authorization", "Bearer my-secret-token")
		token, err := GetBearerToken(headers)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if token != "my-secret-token" {
			t.Fatalf("expected my-secret-token, got %q", token)
		}
	})

	t.Run("extra whitespace", func(t *testing.T) {
		headers := http.Header{}
		headers.Set("Authorization", "Bearer    my-secret-token")
		token, err := GetBearerToken(headers)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if token != "my-secret-token" {
			t.Fatalf("expected my-secret-token, got %q", token)
		}
	})

	t.Run("missing header", func(t *testing.T) {
		headers := http.Header{}
		if _, err := GetBearerToken(headers); err == nil {
			t.Fatal("expected error for missing Authorization header, got nil")
		}
	})

	t.Run("not a bearer token", func(t *testing.T) {
		headers := http.Header{}
		headers.Set("Authorization", "Basic abc123")
		if _, err := GetBearerToken(headers); err == nil {
			t.Fatal("expected error for non-bearer Authorization header, got nil")
		}
	})
}
