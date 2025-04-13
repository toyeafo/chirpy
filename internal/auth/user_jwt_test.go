package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWTAuth(t *testing.T) {
	testID := uuid.New()
	testTokenSecret := "secret"
	testExpiry := 10 * time.Minute
	testToken, err := MakeJWT(testID, testTokenSecret, testExpiry)
	if err != nil {
		t.Fatalf("Failed to create test token: %v", err)
	}

	expiredTokenExpiry := 1 * time.Second
	expiredToken, _ := MakeJWT(testID, testTokenSecret, expiredTokenExpiry)

	time.Sleep(2 * time.Second)

	tests := []struct {
		name    string
		token   string
		secret  string
		wantErr bool
	}{
		{
			name:    "Correct Token",
			token:   testToken,
			secret:  testTokenSecret,
			wantErr: false,
		},
		{
			name:    "Incorrect secret",
			token:   testToken,
			secret:  "psyche",
			wantErr: true,
		},
		{
			name:    "Empty Secret",
			token:   testToken,
			secret:  "",
			wantErr: true,
		},
		{
			name:    "Incorrect Token",
			token:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJrZXkiOiJ2YWx1ZSIsImtleTIiOiJ2YWx1ZTIiLCJpYXQiOjE2MzQxNzgxMTB9.vnXM0oxw05QH1Vs6RsvYp6LaEqFFqZ-NExQMXBgP7Mk",
			secret:  testTokenSecret,
			wantErr: true,
		},
		{
			name:    "Expired Token",
			token:   expiredToken,
			secret:  testTokenSecret,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ValidateJWT(tt.token, tt.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
