package auth

import (
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
