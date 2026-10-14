package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeandValidateJWT(t *testing.T) {
	userID := uuid.New()
	secret := "SECRET"

	token, err := MakeJWT(userID, secret, time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	gotID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("ValidateJWT failed: %v", err)
	}

	if gotID != userID {
		t.Errorf("got user ID %v, want %v", gotID, userID)
	}
}

func TestWrongSecret(t *testing.T) {
	userID := uuid.New()
	correctSecret := "SECRET"
	wrongSecret := "NOT_SECRET"

	token, err := MakeJWT(userID, correctSecret, time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	_, err = ValidateJWT(token, wrongSecret)
	if err == nil {
		t.Error("Expected error for wrong secret but got nil")
	}
}

func TestExpiredToken(t *testing.T) {
	userID := uuid.New()
	secret := "SECRET"

	token, err := MakeJWT(userID, secret, -time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT failed: %v", err)
	}

	_, err = ValidateJWT(token, secret)
	if err == nil {
		t.Error("Expected error for expired token but got nil")
	}
}
