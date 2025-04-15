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

// func TestHashPassword(t *testing.T) {
// 	passwrd := "samuael123"
// 	hash_passwrd, err := HashPassword(passwrd)
// 	if err != nil || hash_passwrd == passwrd {
// 		t.Errorf(`HashPassword(passwrd) = %q %v`, hash_passwrd, err)
// 	}
// }

// func TestCheckPasswordHash(t *testing.T) {
// 	passwrd := "samuel123"
// 	hashPwd, _ := HashPassword(passwrd)
// 	checkPwd := CheckPasswordHash(hashPwd, passwrd)
// 	if checkPwd != nil {
// 		t.Errorf(`CheckPasswordHash(hashPwd, passwrd) = %q should be nil`, checkPwd)
// 	}
// }

func TestCheckPasswordHash(t *testing.T) {
	// First, we need to create some hashed passwords for testing
	password1 := "correctPassword123!"
	password2 := "anotherPassword456!"
	hash1, _ := HashPassword(password1)
	// t.Logf("Generated hash1: %s", hash1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name     string
		password string
		hash     string
		wantErr  bool
	}{
		{
			name:     "Correct password",
			password: password1,
			hash:     hash1,
			wantErr:  false,
		},
		{
			name:     "Incorrect password",
			password: "wrongPassword",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Password doesn't match different hash",
			password: password1,
			hash:     hash2,
			wantErr:  true,
		},
		{
			name:     "Empty password",
			password: "",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Invalid hash",
			password: password1,
			hash:     "invalidhash",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Logf("Testing with hash: %s", tt.hash)
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
