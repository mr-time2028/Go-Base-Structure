package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"testing"
	"time"
)

const (
	tokenExpiry   = 5 * time.Minute
	refreshExpiry = 60 * time.Minute
)

// Helper function to generate an expired token
func generateExpiredToken(secret string) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = "1"
	claims["exp"] = time.Now().Add(-2 * tokenExpiry).Unix()
	signedToken, _ := token.SignedString([]byte(secret))
	return signedToken
}

// Helper function to generate a token signed with a different signing method
func generateTokenWithDifferentSigningMethod() string {
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	token := jwt.New(jwt.SigningMethodRS256) // Use a different signing method (RS256) instead of HMAC
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = "1"
	claims["exp"] = time.Now().Add(tokenExpiry).Unix()
	signedToken, _ := token.SignedString(privateKey)
	return signedToken
}

// Helper function to generate a token with an invalid issuer
func generateTokenWithInvalidIssuer(secret string) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = "1"
	claims["exp"] = time.Now().Add(tokenExpiry).Unix()
	claims["iss"] = "someIssuer"
	signedToken, _ := token.SignedString([]byte(secret))
	return signedToken
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
