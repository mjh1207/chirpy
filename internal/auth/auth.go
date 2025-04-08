package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error){
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func CheckPasswordHash(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func MakeJWT(userID uuid.UUID, tokenSecret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer: "chirpy",
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Hour)),
		Subject: userID.String(),
	})
	return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}
	parsedToken, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})

	if err != nil {
		return uuid.Nil, err
	}
	if !parsedToken.Valid {
		return uuid.Nil, fmt.Errorf("parsed token is not valid")
	}

	userId, err := parsedToken.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	parsedUserId, err := uuid.Parse(userId)
	if err != nil {
		return uuid.Nil, err
	}

	return parsedUserId, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	token := headers.Get("Authorization")
	if token == "" {
		return "", fmt.Errorf("no authorization present")
	}

	return strings.TrimPrefix(token, "Bearer "), nil
}

func MakeRefreshToken() (string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(key), nil
}