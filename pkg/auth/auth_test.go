package auth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestAuth_NewJWTAuth(t *testing.T) {
	jAuth := NewJWTAuth()

	expectedIssuer := "localhost"
	if jAuth.Issuer != expectedIssuer {
		t.Errorf("expected %s as issuer but got %s", expectedIssuer, jAuth.Issuer)
	}

	expectedTokenExpiry := 5 * time.Minute
	if jAuth.TokenExpiry != expectedTokenExpiry {
		t.Errorf("expected %v as token expiry but got %v", expectedTokenExpiry, jAuth.TokenExpiry)
	}
}

func TestAuth_GenerateTokenPair(t *testing.T) {
	auth := NewTestJWTAuth()

	jUser := &JwtUser{
		ID:        1,
		FirstName: "John",
		LastName:  "Smith",
	}

	tokenPairs, err := auth.GenerateTokenPair(jUser)

	if tokenPairs.Token == "" {
		t.Error("access token is empty")
	}
	if tokenPairs.RefreshToken == "" {
		t.Error("refresh token is empty")
	}

	if err != nil {
		t.Errorf("generateTokenPair returned an unexpected error: %v", err)
	}
}

func TestAuth_GetTokenFromHeaderAndVerify(t *testing.T) {
	auth := NewTestJWTAuth()
	jUser := &JwtUser{
		ID:        1,
		FirstName: "John",
		LastName:  "Smith",
	}

	// case 1: valid header
	req, _ := http.NewRequest("POST", "/", nil)
	rr := httptest.NewRecorder()

	tokens, _ := auth.GenerateTokenPair(jUser)

	req.Header.Set("Authorization", "Bearer "+tokens.Token)

	_, _, err := auth.GetTokenFromHeaderAndVerify(rr, req)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	// case 2: invalid auth header
	req, _ = http.NewRequest("POST", "/", nil)
	rr = httptest.NewRecorder()

	req.Header.Set("Authorization", "")

	_, _, err = auth.GetTokenFromHeaderAndVerify(rr, req)
	if err == nil {
		t.Errorf("expected error: there no authorization header, but got no error")
	}

	// case 3: length of auth header is greater than 2
	req, _ = http.NewRequest("POST", "/", nil)
	rr = httptest.NewRecorder()

	req.Header.Set("Authorization", "Bearer some "+tokens.Token)

	_, _, err = auth.GetTokenFromHeaderAndVerify(rr, req)
	if err == nil {
		t.Errorf("expected error: invalid auth header, but got no error")
	}

	// case 4: key "Bearer" does not exist in auth header
	req, _ = http.NewRequest("POST", "/", nil)
	rr = httptest.NewRecorder()

	req.Header.Set("Authorization", "Ber "+tokens.Token)

	_, _, err = auth.GetTokenFromHeaderAndVerify(rr, req)
	if err == nil {
		t.Errorf("expected error: invalid auth header, but got no error")
	}

	// case 5: invalid token (we tested ParseWithClaims already)
	req, _ = http.NewRequest("POST", "/", nil)
	rr = httptest.NewRecorder()

	req.Header.Set("Authorization", "Bearer "+tokens.Token[:10])

	_, _, err = auth.GetTokenFromHeaderAndVerify(rr, req)
	if err == nil {
		t.Errorf("expected an error, but got no error")
	}
}

func TestAuth_ParseWithClaims(t *testing.T) {
	// Set up the Auth instance
	auth := NewTestJWTAuth()

	jUser := &JwtUser{
		ID:        1,
		FirstName: "John",
		LastName:  "Smith",
	}

	// Generate a valid token
	validToken, _ := auth.GenerateTokenPair(jUser)

	// Test case 1: Valid token
	claims, err := auth.ParseWithClaims(validToken.Token)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if claims == nil {
		t.Error("claims should not be nil")
	}

	// Test case 2: Invalid token
	invalidToken := validToken.Token[:10]
	_, err = auth.ParseWithClaims(invalidToken)
	if err == nil {
		t.Error("expected error, but got no error")
	}

	// test case 3: Token signed with an unexpected signing method
	unexpectedSigningMethodToken := generateTokenWithDifferentSigningMethod()
	_, err = auth.ParseWithClaims(unexpectedSigningMethodToken)
	if err == nil || !strings.Contains(err.Error(), "unexpected signing method: RS256") {
		t.Error("expected error: unexpected signing method")
	}

	// test case 4: Expired token
	expiredToken := generateExpiredToken(auth.Secret)
	_, err = auth.ParseWithClaims(expiredToken)
	if err == nil || !strings.Contains(err.Error(), "token has expired") {
		t.Error("expected error: token has expired")
	}

	// test case 5: Invalid issuer
	unexpectedIssuerToken := generateTokenWithInvalidIssuer(auth.Secret)
	_, err = auth.ParseWithClaims(unexpectedIssuerToken)
	if err == nil || !strings.Contains(err.Error(), "invalid issuer") {
		t.Error("expected error: invalid issuer")
	}
}
