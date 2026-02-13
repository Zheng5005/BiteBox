package utils

import (
	"net/http"
	"testing"
)

func TestParseToken_success(t *testing.T) {
	secret := "test-secret"
	expectedID := "user-123"

	tokenStr, err := GenerateMockJWT(expectedID, secret)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	r, _ := http.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "Bearer "+tokenStr)

	userID, err := ParseToken(r, secret)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if userID != expectedID {
		t.Errorf("expected user_id %q, got %q", expectedID, userID)
	}
}

func TestParseToken_missingHeader(t *testing.T) {
	r, _ := http.NewRequest("GET", "/", nil)

	_, err := ParseToken(r, "secret")
	if err == nil {
		t.Fatal("expected error for missing Authorization header")
	}
}

func TestParseToken_wrongSecret(t *testing.T) {
	tokenStr, _ := GenerateMockJWT("user-1", "correct-secret")

	r, _ := http.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "Bearer "+tokenStr)

	_, err := ParseToken(r, "wrong-secret")
	if err == nil {
		t.Fatal("expected error for wrong secret")
	}
}
