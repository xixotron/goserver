package auth_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/xixotron/goserver/internal/auth"
)

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	secret := "secret key"
	validToken, err := auth.MakeJWT(userID, secret, time.Hour)
	if err != nil {
		t.Errorf("Couldn't reate valid token: %v", err)
	}

	expiredToken, err := auth.MakeJWT(userID, secret, time.Nanosecond)
	if err != nil {
		t.Errorf("Couldn't reate valid token: %v", err)
	}
	time.Sleep(time.Millisecond)

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		tokenString string
		tokenSecret string
		wantUserID  uuid.UUID
		wantErr     bool
	}{
		{
			name:        "Valid token returns same user id",
			tokenString: validToken,
			tokenSecret: secret,
			wantUserID:  userID,
			wantErr:     false,
		},
		{
			name:        "Wrong secret results in error",
			tokenString: validToken,
			tokenSecret: "not my secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
		{
			name:        "Invalid token results in error",
			tokenString: "not.A.Valid.Token.String",
			tokenSecret: secret,
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
		{
			name:        "Expired token results in error",
			tokenString: expiredToken,
			tokenSecret: secret,
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := auth.ValidateJWT(tt.tokenString, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error: %v, wantErr: %v", err, tt.wantErr)
				return
			}

			if gotUserID != tt.wantUserID {
				t.Errorf("ValidateJWT() = %v, want %v", gotUserID, tt.wantUserID)
			}
		})
	}
}
