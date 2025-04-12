package auth

import "testing"

func TestHashPassword(t *testing.T) {
	passwrd := "samuael123"
	hash_passwrd, err := HashPassword(passwrd)
	if err != nil || hash_passwrd == passwrd {
		t.Errorf(`HashPassword(passwrd) = %q %v`, hash_passwrd, err)
	}
}

func TestCheckPasswordHash(t *testing.T) {
	passwrd := "samuel123"
	hashPwd, _ := HashPassword(passwrd)
	checkPwd := CheckPasswordHash(hashPwd, passwrd)
	if checkPwd != nil {
		t.Errorf(`CheckPasswordHash(hashPwd, passwrd) = %q should be nil`, checkPwd)
	}
}
