package auth_test

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/xixotron/goserver/internal/auth"
)

func TestJWT(t *testing.T) {
	user1ID := uuid.New()
	secret1 := "hello JWT"
	secret2 := "no more hello programs"

	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		userID         uuid.UUID
		makeSecret     string
		validateSecret string
		wantID         uuid.UUID
		wantExpired    bool
		want           string
		wantErr        bool
	}{
		{
			name:           "using correct secret and without expiration expect same user id",
			userID:         user1ID,
			makeSecret:     secret1,
			validateSecret: secret1,
			wantID:         user1ID,
			wantExpired:    false,
			wantErr:        false,
		},
		{
			name:           "using correct secret and expired token expect expired error",
			userID:         user1ID,
			makeSecret:     secret1,
			validateSecret: secret1,
			wantID:         uuid.Nil,
			wantExpired:    true,
			wantErr:        true,
		},
		{
			name:           "different secret expect error",
			userID:         user1ID,
			makeSecret:     secret1,
			validateSecret: secret2,
			wantID:         uuid.Nil,
			wantExpired:    false,
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			duration := 10 * time.Hour
			if tt.wantExpired {
				duration = 1 * time.Nanosecond
			}
			token, gotErr := auth.MakeJWT(tt.userID, tt.makeSecret, duration)
			if gotErr != nil {
				t.Errorf("MakeJWT() failed: %v", gotErr)
			}

			if tt.wantExpired {
				time.Sleep(10 * time.Millisecond)
			}

			got, gotErr := auth.ValidateJWT(token, tt.validateSecret)
			if gotErr != nil {
				if tt.wantExpired && !strings.Contains(gotErr.Error(), "token is expired") {
					t.Errorf("ValidateJWT() expected expiration got error: %v", gotErr)
				} else if !tt.wantErr {
					t.Errorf("ValidateJWT() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("ValidateJWT() succeeded unexpectedly")
			}

			if tt.wantExpired {
				t.Fatal("token didn't expire and ValidateJWT() succeeded unexpectedly")
			}

			if tt.wantID != got {
				t.Errorf("MakeJWT() = %v, want %v", got, tt.wantID)
			}
		})
	}
}
