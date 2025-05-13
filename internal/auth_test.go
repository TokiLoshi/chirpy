package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	password1 := "correctPassword123!"
	password2 := "anotherPassword456"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name string
		password string 
		hash string 
		wantErr bool
	} {
		{
			name: "Correct password",
			password: password1,
			hash: hash1,
			wantErr: false,
		},
		{
			name: "Incorrect password",
			password: "wrongPassword",
			hash: hash1,
			wantErr: true,
		}, 
		{
			name: "Password doesn't match different hash",
			password: password1, 
			hash: hash2, 
			wantErr: true,
		}, 
		{ 
			name: "Empty password",
			password: "",
			hash: hash1, 
			wantErr: true,
		}, 
		{
			name: "Invalid hash", 
			password: password1, 
			hash: "invalidhash", 
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckPasswordHash(tt.hash, tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMakeAndValidateJWT(t *testing.T) {
	// Setup test data 
	userId := uuid.New()
	secret := "test-secret"
	duration := time.Hour 

	token, err := MakeJWT(userId, secret, duration)
	if err != nil {
		t.Fatalf("Failed to create JWT: %v", err)
	}
	// validate the token 
	returnedId, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("Failed to validate JWT: %v", err)
	}
	if returnedId != userId {
		t.Errorf("Expected user ID %v, got %v", userId, returnedId)
	}
}

func TestExpiredJWT(t *testing.T) {
	userId := uuid.New()
	secret := "test-secret"
	duration := -time.Hour
	token, err := MakeJWT(userId, secret, duration)
	if err != nil {
		t.Fatalf("Failed to make JWT %v", err)
	}
	_, err = ValidateJWT(token, secret)
	if err == nil {
		t.Error("Expected error for failed validationto validate JWT", err)
	}

}

func TestInvalidSecret (t *testing.T) {
	userId := uuid.New()
	correctSecret := "correct-secret"
	wrongSecret := "incorrect-secret"
	duration := time.Hour 
	token, err := MakeJWT(userId, correctSecret, duration)
	if err != nil {
		t.Fatalf("Failed to Make JWT %v", err)
	}

	_, err = ValidateJWT(token, wrongSecret)
	if err == nil {
		t.Error("Expected error for incorrect sercret", err)
	}
}

