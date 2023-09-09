package json

import (
	"encoding/json"
	"errors"
	"go-base-structure/pkg/validators"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestReadJSON(t *testing.T) {
	testCases := []struct {
		name        string
		request     string
		expectedErr bool
	}{
		{
			"valid data",
			`{"first_name": "John", "last_name": "Smith", "age": 21, "is_some": true}`,
			false,
		},
		{
			"invalid data",
			`{"last_name": "Smith", "age": 21, "is_some": true}`,
			true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var requestBody struct {
				FirstName string `json:"first_name" required:"true"`
				LastName  string `json:"last_name"`
				Age       int    `json:"age"`
				IsSome    *bool  `json:"is_some"`
			}

			reqBody := strings.NewReader(tc.request)
			req, _ := http.NewRequest("POST", "/", reqBody)
			rr := httptest.NewRecorder()
			req.Header.Set("Content-Type", "application/json")

			validator := ReadJSON(rr, req, &requestBody)

			if tc.expectedErr && validator.Valid() || !tc.expectedErr && !validator.Valid() {
				t.Errorf("unexpected error: %v", validator.Errors.MessageMap)
			}
		})
	}
}

func TestWriteJSON(t *testing.T) {
	// Create a test data structure
	type TestData struct {
		Message string `json:"message"`
	}
	data := TestData{
		Message: "Hello, World!",
	}

	// Create a test response recorder
	rr := httptest.NewRecorder()

	// Call the WriteJSON function with the test response recorder
	err := WriteJSON(rr, http.StatusOK, data, http.Header{"Custom-Header": []string{"Value"}})
	if err != nil {
		t.Errorf("writeJSON returned an error: %v", err)
	}

	// Validate the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("unexpected status code, got: %d, want: %d", rr.Code, http.StatusOK)
	}

	// Validate the response headers
	expectedContentType := "application/json"
	if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("unexpected Content-Type header, got: %s, want: %s", contentType, expectedContentType)
	}

	expectedCustomHeaderValue := "Value"
	if customHeaderValue := rr.Header().Get("Custom-Header"); customHeaderValue != expectedCustomHeaderValue {
		t.Errorf("unexpected Custom-Header value, got: %s, want: %s", customHeaderValue, expectedCustomHeaderValue)
	}

	// Validate the response body
	var responseBody TestData
	err = json.Unmarshal(rr.Body.Bytes(), &responseBody)
	if err != nil {
		t.Errorf("failed to unmarshal response body: %v", err)
	}

	// Validate the response data
	expectedData := TestData{
		Message: "Hello, World!",
	}
	if !reflect.DeepEqual(responseBody, expectedData) {
		t.Errorf("unexpected response data, got: %+v, want: %+v", responseBody, expectedData)
	}

	// Validate the response body content
	expectedResponseBody := `{"message":"Hello, World!"}`
	if rr.Body.String() != expectedResponseBody {
		t.Errorf("unexpected response body, got: %s, want: %s", rr.Body.String(), expectedResponseBody)
	}

	// Test the error case when marshaling the JSON data
	err = WriteJSON(rr, http.StatusOK, make(chan int))
	if err == nil {
		t.Error("expected an error when marshaling JSON data, but got nil")
	}
}

func TestErrorStrJSON(t *testing.T) {
	rr := httptest.NewRecorder()
	mockError := errors.New("some error")

	// Verify the expected JSON response (default status code is 400)
	expectedResponse := `{"error":true,"message":"some error"}`
	expectedStatusCode := http.StatusBadRequest
	ErrorStrJSON(rr, mockError)

	if rr.Body.String() != expectedResponse {
		t.Errorf("unexpected response body, got: %s, want: %s", rr.Body.String(), expectedResponse)
	}
	if rr.Code != expectedStatusCode {
		t.Errorf("unexpected response code, got %d, want %d", rr.Code, expectedStatusCode)
	}

	// Verify another response status code (change default status code to 500)
	newStatusCode := http.StatusInternalServerError
	ErrorStrJSON(rr, mockError, newStatusCode)

	if rr.Code != expectedStatusCode {
		t.Errorf("unexpected status code, got: %d, want: %d", rr.Code, newStatusCode)
	}
}

func TestErrorMapJSON(t *testing.T) {
	rr := httptest.NewRecorder()
	mockError := map[string][]string{
		"first_name": {
			"this field is required",
		},
		"last_name": {
			"min length error",
			"max length error",
		},
	}

	vError := validators.Errors{
		MessageMap: mockError,
		Code:       http.StatusBadRequest,
	}

	// Verify the expected JSON response (default status code is 400)
	expectedResponse := `{"error":true,"message":{"first_name":["this field is required"],"last_name":["min length error","max length error"]}}`
	expectedStatusCode := http.StatusBadRequest
	ErrorMapJSON(rr, vError)

	if rr.Body.String() != expectedResponse {
		t.Errorf("unexpected response body, got: %s, want: %s", rr.Body.String(), expectedResponse)
	}
	if rr.Code != expectedStatusCode {
		t.Errorf("unexpected response code, got %d, want %d", rr.Code, expectedStatusCode)
	}

	// Verify another response status code (change default status code to 500)
	vError.Code = http.StatusInternalServerError
	ErrorMapJSON(rr, vError)

	if rr.Code != expectedStatusCode {
		t.Errorf("unexpected status code, got: %d, want: %d", rr.Code, vError.Code)
	}
}
