package auth

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"go-base-structure/pkg/env"
	"net/http"
	"strings"
	"time"
)

// Auth contains configs for jwt
type Auth struct {
	Issuer        string
	Audience      string
	Secret        string
	TokenExpiry   time.Duration
	RefreshExpiry time.Duration
}

// JwtUser contains user data for jwt
type JwtUser struct {
	ID        int
	FirstName string
	LastName  string
}

// TokenPairs contains access and refresh tokens
type TokenPairs struct {
	Token        string `json:"access"`
	RefreshToken string `json:"refresh"`
}

// Claims contains default registered claims for jwt
type Claims struct {
	jwt.RegisteredClaims
	TokenType string
}

// NewJWTAuth initial an Auth instance
func NewJWTAuth() *Auth {
	issuer := env.GetEnvOrDefaultString("ISSUER", "localhost")
	audience := env.GetEnvOrDefaultString("AUDIENCE", "localhost")
	secret := env.GetEnvOrDefaultString("SECRET", "")
	tokenExpiry := env.GetEnvOrDefaultInt("TOKEN_EXPIRY", 5)
	refreshExpiry := env.GetEnvOrDefaultInt("REFRESH_EXPIRY", 60)

	return &Auth{
		Issuer:        issuer,
		Audience:      audience,
		Secret:        secret,
		TokenExpiry:   time.Duration(tokenExpiry) * time.Minute,
		RefreshExpiry: time.Duration(refreshExpiry) * time.Minute,
	}
}

// NewTestJWTAuth initial an test auth.Auth instance for test
func NewTestJWTAuth() *Auth {
	return &Auth{
		Issuer:        "test.com",
		Audience:      "test.com",
		Secret:        "testsecret",
		TokenExpiry:   5 * time.Minute,
		RefreshExpiry: 60 * time.Minute,
	}
}

// GenerateTokenPair generates access and refresh tokens
func (a *Auth) GenerateTokenPair(user *JwtUser) (TokenPairs, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	claims["sub"] = fmt.Sprint(user.ID)
	claims["aud"] = a.Audience
	claims["iss"] = a.Issuer
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(a.TokenExpiry).Unix()
	claims["TokenType"] = "access"

	signedToken, err := token.SignedString([]byte(a.Secret))
	if err != nil {
		return TokenPairs{}, err
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshTokenClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshTokenClaims["sub"] = fmt.Sprint(user.ID)
	refreshTokenClaims["iss"] = a.Issuer
	refreshTokenClaims["iat"] = time.Now().Unix()
	refreshTokenClaims["exp"] = time.Now().Add(a.RefreshExpiry).Unix()
	refreshTokenClaims["TokenType"] = "refresh"

	signedRefreshToken, err := refreshToken.SignedString([]byte(a.Secret))
	if err != nil {
		return TokenPairs{}, err
	}

	var tokenPairs = TokenPairs{
		Token:        signedToken,
		RefreshToken: signedRefreshToken,
	}

	return tokenPairs, nil
}

// GetTokenFromHeaderAndVerify extract access token from client request and validate it
func (a *Auth) GetTokenFromHeaderAndVerify(w http.ResponseWriter, r *http.Request) (string, *Claims, error) {
	w.Header().Add("Vary", "Authorization")

	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		return "", nil, errors.New("there no authorization header")
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		return "", nil, errors.New("invalid auth header")
	}

	if headerParts[0] != "Bearer" {
		return "", nil, errors.New("invalid auth header")
	}

	token := headerParts[1]

	claims, err := a.ParseWithClaims(token)
	if err != nil {
		return "", nil, err
	}

	return token, claims, nil
}

// ParseWithClaims parse jwt using secret key
func (a *Auth) ParseWithClaims(token string) (*Claims, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(a.Secret), nil
	})
	if err != nil {
		if strings.Contains(err.Error(), "token is expired") {
			return nil, errors.New("token has expired")
		}
		return nil, err
	}

	if claims.Issuer != a.Issuer {
		return nil, errors.New("invalid issuer")
	}

	return claims, nil
}
