package auth

import (
	"net/http"
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

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name      string
		headers   http.Header
		wantToken string
		wantErr   bool
	}{
		{
			name: "Valid Bearer token",
			headers: http.Header{
				"Authorization": []string{"Bearer valid_token"},
			},
			wantToken: "valid_token",
			wantErr:   false,
		},
		{
			name:      "Missing Authorization Header",
			headers:   http.Header{},
			wantToken: "",
			wantErr:   true,
		},
		{
			name: "Malformed Authorization Header",
			headers: http.Header{
				"Authorization": []string{"InvalidBearer token"},
			},
			wantToken: "",
			wantErr:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotToken, err := GetBearerToken(test.headers)
			if (err != nil) != test.wantErr {
				t.Errorf("getBearerToken() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if gotToken != test.wantToken {
				t.Errorf("GetBearerToken() gotToken = %v, want %v", gotToken, test.wantToken)
			}
		})
	}
}
