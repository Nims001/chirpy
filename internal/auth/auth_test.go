package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

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
