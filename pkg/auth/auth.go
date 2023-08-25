package auth

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
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
}

// GenerateTokenPair generates access and refresh tokens
func (a *Auth) GenerateTokenPair(user *JwtUser) (TokenPairs, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	claims["sub"] = fmt.Sprint(user.ID)
	claims["aud"] = a.Audience
	claims["iss"] = a.Issuer
	claims["iat"] = time.Now().UTC().Unix()
	claims["typ"] = "JWT"
	claims["exp"] = time.Now().UTC().Add(a.TokenExpiry).Unix()

	signedToken, err := token.SignedString([]byte(a.Secret))
	if err != nil {
		return TokenPairs{}, err
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshTokenClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshTokenClaims["sub"] = fmt.Sprint(user.ID)
	refreshTokenClaims["iss"] = a.Issuer
	refreshTokenClaims["iat"] = time.Now().UTC().Unix()
	refreshTokenClaims["exp"] = time.Now().UTC().Add(a.RefreshExpiry).Unix()

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

	headerParts := strings.Split(authHeader, "")
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

	if claims.Issuer != a.Issuer {
		return "", nil, errors.New("invalid issuer")
	}

	if claims.ExpiresAt.Before(time.Now()) {
		return "", nil, errors.New("token has expired")

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
		if strings.HasPrefix(err.Error(), "token is expired by") {
			return nil, errors.New("token has expired")
		}
		return nil, err
	}
	return claims, nil
}
