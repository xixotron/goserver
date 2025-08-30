package auth_test

import (
	"testing"

	"github.com/xixotron/goserver/internal/auth"
)

func TestHashPassword(t *testing.T) {
	type test struct {
		name string // description of this test case
		// Named input parameters for target function.
		password1 string
		password2 string
		wantErr1  bool
		wantErr2  bool
	}

	tests := []test{
		{
			name:      "check empty password against empty password",
			password1: "",
			password2: "",
			wantErr1:  false,
			wantErr2:  false,
		}, {
			name:      "check short password against itself",
			password1: "hello",
			password2: "hello",
			wantErr1:  false,
			wantErr2:  false,
		}, {
			name:      "check spaces in password against itself",
			password1: "hello world",
			password2: "hello world",
			wantErr1:  false,
			wantErr2:  false,
		}, {
			name:      "check spaces in password against empty",
			password1: "hello world",
			password2: "",
			wantErr1:  false,
			wantErr2:  true,
		}, {
			name:      "check short password against empty",
			password1: "hello",
			password2: "",
			wantErr1:  false,
			wantErr2:  true,
		}, {
			name:      "check short password against other short password",
			password1: "hello",
			password2: "world",
			wantErr1:  false,
			wantErr2:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, gotErr1 := auth.HashPassword(tt.password1)
			if gotErr1 != nil {
				if !tt.wantErr1 {
					t.Errorf("HashPassword() failed: %v", gotErr1)
				}
				return
			}
			if tt.wantErr1 {
				t.Fatal("HashPassword() succeeded unexpectedly")
			}

			gotErr2 := auth.CheckPasswordHash(tt.password2, hash)
			if gotErr2 != nil {
				if !tt.wantErr2 {
					t.Errorf("CheckPasswordHash() failed: %v", gotErr2)
				}
				return
			}
			if tt.wantErr2 {
				t.Fatal("CheckPasswordHash() succeeded unexpectedly")
			}
		})
	}
}
