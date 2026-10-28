package auth

import (
	"net/http"
	"testing"
)

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
