package auth

//package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// HashPassword -
func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hash, nil
}

// CheckPasswordHash -
func CheckPasswordHash(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}
	return match, nil
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {

	// Create claims with multiple fields populated
	claims := &jwt.RegisteredClaims{
		// A usual scenario is to set the expiration time relative to the current time
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "chirpy-access",
		Subject:   fmt.Sprint(userID),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims) //new with claims returns a token type struct
	ss, err := token.SignedString([]byte(tokenSecret))
	return ss, err

}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {

	// the parse function takes the token string and a callback function that returns the key for validation.
	// ParseWithClaims — the parser calls your function, giving it the parsed token, and your function returns the key that should be used to check the signature.
	// the empty claims struct is passed to the parse function, which will populate it with the claims from the token if the token is valid.
	parsedToken, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	if claims, ok := parsedToken.Claims.(*jwt.RegisteredClaims); ok && parsedToken.Valid {
		userID, err := uuid.Parse(claims.Subject)
		if err != nil {
			return uuid.Nil, err
		}
		return userID, nil
	} else {
		return uuid.Nil, fmt.Errorf("invalid token")
	}
}

func TestToken(t *testing.T) {
	tokenSecret := "mysecretkey"
	userID := uuid.New()
	// Generate a token
	token, err := MakeJWT(userID, tokenSecret, time.Hour)
	if err != nil {
		t.Errorf("Failed to generate token: %v", err)
	} else {
		t.Logf("Generated token: %s", token)
	}

	// Validate the token
	validatedUserID, err := ValidateJWT(token, tokenSecret)
	if err != nil {
		t.Errorf("Failed to validate token: %v", err)
	}

	if validatedUserID != userID {
		t.Errorf("Expected userID %v, got %v", userID, validatedUserID)
	} else {
		t.Logf("Successfully validated token for userID: %v", validatedUserID)
	}

	expiredToken, err := MakeJWT(userID, tokenSecret, -time.Hour) // Create an expired token
	if err != nil {
		t.Errorf("Failed to generate expired token: %v", err)
	}

	_, err = ValidateJWT(expiredToken, tokenSecret)
	if err == nil {
		t.Errorf("Expected error for expired token, got nil")
	} else {
		t.Logf("Successfully caught expired token: %v", err)
	}

	_, err = ValidateJWT(token, "wrongsecret")
	if err == nil {
		t.Errorf("Expected error for invalid token")
	} else {
		t.Logf("Successfully caught invalid token: %v", err)
	}
}
