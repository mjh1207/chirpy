package auth

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

// Test MakeJWT function
func TestJWTCreationAndValidation(t *testing.T) {
	// Generate a random UUID
	userID := uuid.New()
	secret := "test-secret"

	// Create a token
	token, err := MakeJWT(userID, secret, time.Hour)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	// Validate the token
	returnedUserID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}
	
	// Check if returned user ID matches the original
	if returnedUserID != userID {
		t.Errorf("Expected user ID %s, got %s", userID, returnedUserID)
	}
}

func TestExpiredToken(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"

	// Create a token that expires after one second
	token, err := MakeJWT(userID, secret, time.Second)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	// Wait for token to expire
	time.Sleep(2 * time.Second)

	// Attempt to validate expired token
	_, err = ValidateJWT(token, secret)
	if err == nil {
		t.Fatalf("Expected token to be invalid, got VALID")
	}
}

func TestInvalidSecret(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"
	invalidSecret := "invalid-secret"
	
	// Create a token
	token, err := MakeJWT(userID, secret, time.Hour)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}

	// Attempt to validate token with incorrect secret
	_, err = ValidateJWT(token, invalidSecret)
	if err == nil {
		t.Fatalf("Expected error when validating token, got nil")
	}
}

func TestBearerToken(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"

	token, err := MakeJWT(userID, secret, time.Second)
	if err != nil {
		t.Fatalf("Failed to create token: %v", err)
	}


	req, err := http.NewRequest("GET", "https://example.com", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	token, err = GetBearerToken(req.Header)
	if err != nil {
		t.Fatalf("Failed to get authorization from headers: %v", err)
	}

	_, err = ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("Failed to validate bearer token")
	}
}