package auth

import "testing"

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
