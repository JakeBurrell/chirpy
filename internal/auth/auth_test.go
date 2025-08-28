package auth

import (
	"github.com/google/uuid"
	"net/http"
	"testing"
	"time"
)

func TestCheckPasswordHash(t *testing.T) {
	password1 := "CorrectPassword1!"
	password2 := "WrongPassword1!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)
	tests := []struct {
		name     string
		password string
		hash     string
		wantErr  bool
	}{
		{
			name:     "Correct Password",
			password: password1,
			hash:     hash1,
			wantErr:  false,
		},
		{
			name:     "Incorrect Password",
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
			name:     "Second correct password",
			password: password2,
			hash:     hash2,
			wantErr:  false,
		},
		{
			name:     "Invalid hash",
			password: password1,
			hash:     "lol",
			wantErr:  true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := CheckPasswordHash(test.password, test.hash)
			if (err != nil) != test.wantErr {
				t.Errorf("CheckPasswordHash() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	userId := uuid.New()
	validToken, _ := MakeJWT(userId, "secret", time.Hour)
	expiredToken, _ := MakeJWT(userId, "secret", -time.Hour)

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
			tokenSecret: "secret",
			wantUserId:  userId,
			wantErr:     false,
		},
		{
			name:        "Invalid Token",
			tokenString: "Random token",
			tokenSecret: "secret",
			wantUserId:  uuid.Nil,
			wantErr:     true,
		},
		{
			name:        "Wrong secret",
			tokenString: validToken,
			tokenSecret: "wrong",
			wantUserId:  uuid.Nil,
			wantErr:     true,
		},
		{
			name:        "Expired Token",
			tokenString: expiredToken,
			tokenSecret: "secret",
			wantUserId:  uuid.Nil,
			wantErr:     true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotUserId, err := ValidateJWT(test.tokenString, test.tokenSecret)
			if (err != nil) != test.wantErr {
				t.Errorf("ValidateJWT(), error =  %v, wantErr %v", err, test.wantErr)
				return
			}
			if gotUserId != test.wantUserId {
				t.Errorf("ValidateJWT(), gotUerID = %v, want %v", gotUserId, test.wantUserId)
			}
		})
	}
}

func TestBearerToken(t *testing.T) {

	testToken := "testtoken"
	validHeader := http.Header{}
	validHeader.Set("Authorization", "Bearer "+testToken)
	noAuthHeader := http.Header{}
	malformedHeader := http.Header{}
	malformedHeader.Set("Authorization", testToken)

	tests := []struct {
		name     string
		header   http.Header
		expected string
		wantErr  bool
	}{
		{
			name:     "Test Valid Header",
			header:   validHeader,
			expected: testToken,
			wantErr:  false,
		},
		{
			name:     "Test header with no Authorization",
			header:   noAuthHeader,
			expected: "",
			wantErr:  true,
		},
		{
			name:     "Malformed Authorization Header",
			header:   malformedHeader,
			expected: "",
			wantErr:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			token, err := GetBearerToken(test.header)
			if (err != nil) != test.wantErr {
				t.Errorf("GetBearerToken(), error = %v, wantErr = %v", err, test.wantErr)
				return
			}
			if token != test.expected {
				t.Errorf("GetBearerToken(), token = %v want = %v", token, test.expected)
			}

		})
	}
}
