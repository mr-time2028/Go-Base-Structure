package user

import (
	"encoding/json"
	"fmt"
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

func TestLogin(t *testing.T) {
	var theTests = []struct {
		name               string
		requestBody        string
		expectedStatusCode int
	}{
		{"valid data", `{"email": "hamid@test.com", "password": "testPass"}`, http.StatusOK},
		{"no user row", `{"email": "norows@test.com", "password": "testPass"}`, http.StatusUnauthorized},
		{"no body", ``, http.StatusBadRequest},
		{"no password", `{"email": "any@example.com"}`, http.StatusBadRequest},
		{"no email", `{"password": "randompassword"}`, http.StatusBadRequest},
		{"bad json", `{"password: "randompassword"}`, http.StatusBadRequest},
	}

	for _, e := range theTests {
		rr := FakePostRequest(e.requestBody, "/login", http.HandlerFunc(Login))
		if e.expectedStatusCode != rr.Code {
			t.Errorf("%s: returned wrong status code; expected %d but got %d", e.name, e.expectedStatusCode, rr.Code)
		}
	}
}

func TestRefreshToken(t *testing.T) {
	// request to login to receive access and refresh token
	rr := FakePostRequest(`{"email": "hamid@test.com", "password": "testPass"}`, "/login", http.HandlerFunc(Login))
	var response struct {
		RefreshToken string `json:"refresh"`
		AccessToken  string `json:"access"`
	}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("error parsing response: %v", err)
	}

	// case 1: test with valid refresh token
	rr = FakePostRequest(fmt.Sprintf(`{"refresh": "%s"}`, response.RefreshToken), "/refresh", http.HandlerFunc(RefreshToken))
	if http.StatusOK != rr.Code {
		t.Errorf("%s: returned wrong status code; expected %d but got %d", "valid data", http.StatusOK, rr.Code)
	}

	// case 2: send access instead of refresh
	rr = FakePostRequest(fmt.Sprintf(`{"refresh": "%s"}`, response.AccessToken), "/refresh", http.HandlerFunc(RefreshToken))
	if http.StatusUnauthorized != rr.Code {
		t.Errorf("%s: returned wrong status code; expected %d but got %d", "access instead refresh", http.StatusUnauthorized, rr.Code)
	}
}
