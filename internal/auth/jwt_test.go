package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestValidateJwt(t *testing.T) {
	userId := uuid.New()
	validToken, _ := MakeJwt(userId, "secret", time.Minute)

	tests := []struct {
		name        string
		tokenString string
		tokenSecret string
		wantUserId  uuid.UUID
		wantErr     bool
	}{
		{
			name:        "Valid Token",
			tokenString: validToken,
			wantUserId:  userId,
			tokenSecret: "secret",
			wantErr:     false,
		},
		{
			name:        "InValid Token",
			tokenString: "invalid-token-string",
			wantUserId:  uuid.Nil,
			tokenSecret: "secret",
			wantErr:     true,
		},
		{
			name:        "Invalid secret",
			tokenString: validToken,
			wantUserId:  uuid.Nil,
			tokenSecret: "wrong-secret",
			wantErr:     true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotUserId, err := ValidateJwt(tc.tokenString, tc.tokenSecret)
			if (err != nil) != tc.wantErr {
				t.Errorf("ValidateJwt() error = %v, wantError %v", err, tc.wantErr)
				return
			}
			if gotUserId != tc.wantUserId {
				t.Errorf("ValidateJwt() gotUserId = %v, wantUserId =%v", gotUserId, tc.wantUserId)
			}
		})
	}

}

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name      string
		header    http.Header
		wantToken string
		wantErr   bool
	}{
		{
			name:      "valid bearer token",
			header:    http.Header{"Authorization": []string{"Bearer valid-token"}},
			wantToken: "valid-token",
			wantErr:   false,
		},
		{
			name:      "missing authorization header",
			header:    http.Header{},
			wantToken: "",
			wantErr:   true,
		},
		{
			name:      "invalid bearer format",
			header:    http.Header{"Authorization": []string{"invalid-token"}},
			wantToken: "",
			wantErr:   true,
		},
		{
			name:      "empty authorization header",
			header:    http.Header{"Authorization": []string{""}},
			wantToken: "",
			wantErr:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotToken, err := GetBearerToken(tc.header)
			if (err != nil) != tc.wantErr {
				t.Errorf("GetBearerToken() error = %v, wantErr %v", err, tc.wantErr)
			}
			if gotToken != tc.wantToken {
				t.Errorf("GetBearerToken() gotToken = %v, wantToken = %v", gotToken, tc.wantToken)
			}
		})
	}
}
