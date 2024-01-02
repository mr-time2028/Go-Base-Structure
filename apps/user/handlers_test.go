package user

import (
	"encoding/json"
	"fmt"
	"go-base-structure/pkg/auth"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func FakePostRequest(body string, url string, handler http.Handler) *httptest.ResponseRecorder {
	reader := strings.NewReader(body)
	req, _ := http.NewRequest("POST", url, reader)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	return rr
}

func TestRegister(t *testing.T) {
	var testCases = []struct {
		name               string
		requestBody        string
		expectedStatusCode int
		expectedErr        bool
	}{
		{
			"valid data",
			`{
				"email": "ramin@test.com",
				"first_name": "Ramin",
				"last_name": "Djawadi",
				"password": "RaminPass1234",
				"confirm_password": "RaminPass1234"
			}`,
			http.StatusOK,
			false,
		},
		{
			"duplicate email",
			`{
				"email": "John@test.com",
				"first_name": "John",
				"last_name": "Smith",
				"password": "John1234",
				"confirm_password": "John1234"
			}`,
			http.StatusBadRequest,
			true,
		},
		{
			"passwords are not match",
			`{
				"email": "ramin@test.com",
				"first_name": "Ramin",
				"last_name": "Djawadi",
				"password": "RaminPass1234",
				"confirm_password": "Pass1234"
			}`,
			http.StatusBadRequest,
			true,
		},
		{
			"not all fields in json",
			`{
				"email": "ramin@test.com",
			}`,
			http.StatusBadRequest,
			true,
		},
		{
			"invalid email address",
			`{
				"email": "ramin",
				"first_name": "Ramin",
				"last_name": "Djawadi",
				"password": "RaminPass1234",
				"confirm_password": "RaminPass1234"
			}`,
			http.StatusBadRequest,
			true,
		},
	}

	for _, e := range testCases {
		rr := FakePostRequest(e.requestBody, "/register", http.HandlerFunc(Register))
		if e.expectedStatusCode != rr.Code {
			t.Errorf("%s: returned wrong status code; expected %d but got %d, %v", e.name, e.expectedStatusCode, rr.Code, rr.Body)
		}
	}
}

func TestLogin(t *testing.T) {
	var theTests = []struct {
		name               string
		requestBody        string
		expectedStatusCode int
	}{
		{
			"valid data",
			`{"email": "admin@test.com", "password": "FAdminPass"}`,
			http.StatusOK,
		},
		{
			"no user row",
			`{"email": "norows@test.com", "password": "testPass"}`,
			http.StatusUnauthorized,
		},
		{
			"wrong password",
			`{"email": "David@test.com", "password": "WrongPass"}`,
			http.StatusUnauthorized,
		},
		{
			"no body",
			``,
			http.StatusBadRequest,
		},
		{
			"no password",
			`{"email": "any@example.com"}`,
			http.StatusBadRequest,
		},
		{
			"no email",
			`{"password": "randomPass"}`,
			http.StatusBadRequest,
		},
		{
			"bad json",
			`{"password: "randomPass"}`,
			http.StatusBadRequest,
		},
	}

	for _, e := range theTests {
		rr := FakePostRequest(e.requestBody, "/login", http.HandlerFunc(Login))
		if e.expectedStatusCode != rr.Code {
			t.Errorf("%s: returned wrong status code; expected %d but got %d", e.name, e.expectedStatusCode, rr.Code)
		}
	}
}

func TestRefreshToken(t *testing.T) {
	// request to login
	var response struct {
		RefreshToken string `json:"refresh"`
		AccessToken  string `json:"access"`
	}
	rr := FakePostRequest(`{"email": "admin@test.com", "password": "FAdminPass"}`, "/login", http.HandlerFunc(Login))
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("error parsing response: %v", err)
	}

	// test refresh path
	var theTests = []struct {
		name               string
		requestBody        string
		expectedStatusCode int
	}{
		{
			"valid refresh",
			fmt.Sprintf(`{"refresh": "%s"}`, response.RefreshToken),
			http.StatusOK,
		},
		{
			"access instead of refresh",
			fmt.Sprintf(`{"refresh": "%s"}`, response.AccessToken),
			http.StatusUnauthorized,
		},
		{
			"invalid refresh",
			fmt.Sprintf(`{"refresh": "%s"}`, response.RefreshToken[:10]),
			http.StatusUnauthorized,
		},
		{
			"no body",
			``,
			http.StatusBadRequest,
		},
		{
			"bad json",
			`{"password: "randomPass"}`,
			http.StatusBadRequest,
		},
	}

	for _, e := range theTests {
		rr = FakePostRequest(e.requestBody, "/refresh", http.HandlerFunc(RefreshToken))
		if e.expectedStatusCode != rr.Code {
			t.Errorf("%s: returned wrong status code; expected %d but got %d, err is", e.name, e.expectedStatusCode, rr.Code)
		}
	}

	// test when user not in database but has a valid refresh token (id 240 is not in database)
	jUser := &auth.JwtUser{
		ID:        240,
		FirstName: "Alex",
		LastName:  "Parker",
	}
	tokens, _ := userApp.Auth.GenerateTokenPair(jUser)
	rr = FakePostRequest(fmt.Sprintf(`{"refresh": "%s"}`, tokens.RefreshToken), "/refresh", http.HandlerFunc(RefreshToken))
	if rr.Code != http.StatusUnauthorized {
		t.Errorf("%s: returned wrong status code; expected %d but got %d", "valid refresh, invalid user", http.StatusUnauthorized, rr.Code)
	}
}
