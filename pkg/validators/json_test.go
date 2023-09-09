package validators

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestValidation_MinLengthTag(t *testing.T) {
	var testCases = []struct {
		name           string
		request        string
		expectedErr    bool
		expectedErrMap map[string][]string
	}{
		{
			"valid data",
			`{"first_name": "Benjamin", "last_name": "Montgomery", "age": 21}`,
			false,
			map[string][]string{},
		},
		{
			"invalid minimum character for first_name",
			`{"first_name": "John", "last_name": "Montgomery", "age": 54}`,
			true,
			map[string][]string{"first_name": {fmt.Sprintf(MinErrMsg, 6)}},
		},
		{
			"invalid minimum character for first_name and last_name",
			`{"first_name": "John", "last_name": "Smith", "age": 54}`,
			true,
			map[string][]string{
				"first_name": {fmt.Sprintf(MinErrMsg, 6)},
				"last_name":  {fmt.Sprintf(MinErrMsg, 7)},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var testBody struct {
				FirstName string `json:"first_name" min:"6"`
				LastName  string `json:"last_name" min:"7"`
				City      string `json:"city"`
				Age       int    `json:"age"`
			}

			req, _ := http.NewRequest("POST", "/", strings.NewReader(tc.request))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			maxBytes := 1048576 // one megabyte
			req.Body = http.MaxBytesReader(rr, req.Body, int64(maxBytes))

			dec := json.NewDecoder(req.Body)
			dec.DisallowUnknownFields()
			_ = dec.Decode(&testBody)

			dataType := reflect.TypeOf(&testBody).Elem()
			dataValue := reflect.ValueOf(&testBody).Elem()

			validator := New()
			for i := 0; i < dataType.NumField(); i++ {
				field := dataType.Field(i)
				fieldValue := dataValue.Field(i)

				validator.MinLengthTag(field, fieldValue)
			}

			if tc.expectedErr && validator.Valid() || !tc.expectedErr && !validator.Valid() {
				t.Errorf("unexpected error: %v", validator.Errors.MessageMap)
			} else if len(tc.expectedErrMap) != len(validator.Errors.MessageMap) {
				t.Errorf("unexpected error, want %v, but got %v", tc.expectedErrMap, validator.Errors.MessageMap)
			} else if !reflect.DeepEqual(tc.expectedErrMap, validator.Errors.MessageMap) {
				t.Errorf("expected error %s, but got %s", tc.expectedErrMap, validator.Errors.MessageMap)
			}
		})
	}
}

func TestValidation_MaxLengthTag(t *testing.T) {
	var testCases = []struct {
		name           string
		request        string
		expectedErr    bool
		expectedErrMap map[string][]string
	}{
		{
			"valid data",
			`{"first_name": "John", "last_name": "Smith", "age": 21}`,
			false,
			map[string][]string{},
		},
		{
			"invalid maximum character for first_name",
			`{"first_name": "Jonathan", "last_name": "Smith", "age": 54}`,
			true,
			map[string][]string{"first_name": {fmt.Sprintf(MaxErrMsg, 6)}},
		},
		{
			"invalid maximum character for first_name and last_name",
			`{"first_name": "Benjamin", "last_name": "Montgomery", "age": 54}`,
			true,
			map[string][]string{
				"first_name": {fmt.Sprintf(MaxErrMsg, 6)},
				"last_name":  {fmt.Sprintf(MaxErrMsg, 5)},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var testBody struct {
				FirstName string `json:"first_name" max:"6"`
				LastName  string `json:"last_name" max:"5"`
				City      string `json:"city"`
				Age       int    `json:"age"`
			}

			req, _ := http.NewRequest("POST", "/", strings.NewReader(tc.request))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			maxBytes := 1048576 // one megabyte
			req.Body = http.MaxBytesReader(rr, req.Body, int64(maxBytes))

			dec := json.NewDecoder(req.Body)
			dec.DisallowUnknownFields()
			_ = dec.Decode(&testBody)

			dataType := reflect.TypeOf(&testBody).Elem()
			dataValue := reflect.ValueOf(&testBody).Elem()

			validator := New()
			for i := 0; i < dataType.NumField(); i++ {
				field := dataType.Field(i)
				fieldValue := dataValue.Field(i)

				validator.MaxLengthTag(field, fieldValue)
			}

			if tc.expectedErr && validator.Valid() || !tc.expectedErr && !validator.Valid() {
				t.Errorf("unexpected error: %v", validator.Errors.MessageMap)
			} else if len(tc.expectedErrMap) != len(validator.Errors.MessageMap) {
				t.Errorf("unexpected error, want %v, but got %v", tc.expectedErrMap, validator.Errors.MessageMap)
			} else if !reflect.DeepEqual(tc.expectedErrMap, validator.Errors.MessageMap) {
				t.Errorf("expected error %s, but got %s", tc.expectedErrMap, validator.Errors.MessageMap)
			}
		})
	}
}

func TestValidation_RequiredTag(t *testing.T) {
	var testCases = []struct {
		name           string
		request        string
		expectedErr    bool
		expectedErrMap map[string][]string
	}{
		{
			"valid data",
			`{"first_name": "Benjamin", "last_name": "Montgomery"}`,
			false,
			map[string][]string{},
		},
		{
			"no last name",
			`{"first_name": "John"}`,
			true,
			map[string][]string{"last_name": {RequiredErrMsg}},
		},
		{
			"no first name",
			`{"last_name": "Smith"}`,
			true,
			map[string][]string{"first_name": {RequiredErrMsg}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var testBody struct {
				FirstName string `json:"first_name" required:"true"`
				LastName  string `json:"last_name" required:"true"`
				City      string `json:"city"`
				Age       int    `json:"age"`
			}

			req, _ := http.NewRequest("POST", "/", strings.NewReader(tc.request))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			maxBytes := 1048576 // one megabyte
			req.Body = http.MaxBytesReader(rr, req.Body, int64(maxBytes))

			dec := json.NewDecoder(req.Body)
			dec.DisallowUnknownFields()
			_ = dec.Decode(&testBody)

			dataType := reflect.TypeOf(&testBody).Elem()
			dataValue := reflect.ValueOf(&testBody).Elem()

			validator := New()
			for i := 0; i < dataType.NumField(); i++ {
				field := dataType.Field(i)
				fieldValue := dataValue.Field(i)

				validator.RequiredTag(field, fieldValue)
			}

			if tc.expectedErr && validator.Valid() || !tc.expectedErr && !validator.Valid() {
				t.Errorf("unexpected error: %v", validator.Errors.MessageMap)
			} else if len(tc.expectedErrMap) != len(validator.Errors.MessageMap) {
				t.Errorf("unexpected error, want %v, but got %v", tc.expectedErrMap, validator.Errors.MessageMap)
			} else if !reflect.DeepEqual(tc.expectedErrMap, validator.Errors.MessageMap) {
				t.Errorf("expected error %s, but got %s", tc.expectedErrMap, validator.Errors.MessageMap)
			}
		})
	}
}

func TestValidation_JsonValidation(t *testing.T) {
	testCases := []struct {
		name           string
		request        string
		expectedErr    bool
		expectedErrMap map[string][]string
	}{
		{
			"valid json",
			`{"id": "123456789", "first_name": "John", "last_name": "Smith", "age": 30, "is_some": false}`,
			false,
			map[string][]string{},
		},
		// json validation
		{
			"more than one json",
			`{"age": 30}, {"age": 45}`,
			true,
			map[string][]string{
				"json": {JSONValueErrMsg},
			},
		},
		{
			"invalid json syntax",
			`{"name": "John", age: 30}`,
			true,
			map[string][]string{
				"json": {fmt.Sprintf(JSONSyntaxErrMsg, strconv.Itoa(18))},
			},
		},
		{
			// added some logic based on test name
			"nil request body",
			"",
			true,
			map[string][]string{
				"json": {JSONEmptyBodyErrMsg},
			},
		},
		{
			// added some logic based on test name
			"nil data",
			"",
			true,
			map[string][]string{
				"json": {JSONEmptyBodyErrMsg},
			},
		},
		{
			"cannot decode json",
			`{`,
			true,
			map[string][]string{
				"json": {"unexpected EOF"},
			},
		},
		{
			"invalid data type (age as string)",
			`{"last_name": "Smith", "age": "30", "is_some": false}`,
			true,
			map[string][]string{
				"json": {fmt.Sprintf(JSONUnmarshalTypeErrMsg, strconv.Itoa(34))},
			},
		},
		{
			"unknown field",
			`{"city": "Tokyo"}`,
			true,
			map[string][]string{
				"json": {fmt.Sprintf(JSONUnknownField, `"city"`)},
			},
		},
		// required validation
		{
			"required age (int type) not in json",
			`{"id": "123456789", "first_name": "John", "last_name": "Smith", "is_some": false}`,
			true,
			map[string][]string{
				"age": {RequiredErrMsg},
			},
		},
		{
			`required last_name (string type) have "" value ("" is zero value for string)`,
			`{"id": "123456789", "first_name": "John", "last_name": "", "age": 30, "is_some": false}`,
			true,
			map[string][]string{
				"last_name": {RequiredErrMsg},
			},
		},
		{
			"required age (int type) and last_name (string type) not in json",
			`{"id": "123456789", "first_name": "John", "is_some": false}`,
			true,
			map[string][]string{
				"last_name": {RequiredErrMsg},
				"age":       {RequiredErrMsg},
			},
		},
		{
			"first_name and is_some are not required and not in json",
			`{"id": "123456789", "last_name": "Smith", "age": 21}`,
			false,
			map[string][]string{},
		},
		{
			`first_name have a "" value (zero value for string)`,
			`{"id": "123456789", "first_name": "", "last_name": "Smith", "age": 21}`,
			false,
			map[string][]string{},
		},

		// minimum validation
		{
			"minimum last_name in json",
			`{"id": "123456789", "first_name": "John", "last_name": "Park", "age": 30, "is_some": false}`,
			false,
			map[string][]string{},
		},
		{
			"minimum last_name not in json",
			`{"id": "123456789", "first_name": "John", "last_name": "Lee", "age": 30, "is_some": false}`,
			true,
			map[string][]string{
				"last_name": {
					fmt.Sprintf(MinErrMsg, 4),
				},
			},
		},
		// maximum validation
		{
			"maximum first_name in in json",
			`{"id": "123456789", "first_name": "John", "last_name": "Smith", "age": 30, "is_some": false}`,
			false,
			map[string][]string{},
		},
		{
			"maximum first_name in json but last_name not",
			`{"id": "123456789", "first_name": "John", "last_name": "Montgomery", "age": 30, "is_some": false}`,
			true,
			map[string][]string{
				"last_name": {
					fmt.Sprintf(MaxErrMsg, 8),
				},
			},
		},
		{
			"maximum first_name and last_name not in json",
			`{"id": "123456789", "first_name": "Benjamin", "last_name": "Montgomery", "age": 30, "is_some": false}`,
			true,
			map[string][]string{
				"first_name": {fmt.Sprintf(MaxErrMsg, 7)},
				"last_name":  {fmt.Sprintf(MaxErrMsg, 8)},
			},
		},
		{
			"maximum id and first_name and last_name not in json",
			`{"id": "00123456789", "first_name": "Benjamin", "last_name": "Montgomery", "age": 30, "is_some": false}`,
			true,
			map[string][]string{
				"id":         {fmt.Sprintf(MaxErrMsg, 10)},
				"first_name": {fmt.Sprintf(MaxErrMsg, 7)},
				"last_name":  {fmt.Sprintf(MaxErrMsg, 8)},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var testBody struct {
				ID        string `json:"id" required:"true" max:"10"`
				FirstName string `json:"first_name" max:"7"`
				LastName  string `json:"last_name" required:"true" min:"4" max:"8"`
				Age       int    `json:"age" required:"true"`
				IsSome    *bool  `json:"is_some"`
			}

			reqBody := strings.NewReader(tc.request)
			req, _ := http.NewRequest("POST", "/", reqBody)
			req.Header.Set("Content-Type", "application/json")

			var request interface{}
			request = &testBody
			if tc.name == "nil request body" {
				req.Body = nil
			}
			if tc.name == "nil data" {
				request = nil
			}

			validator := New()
			validator.JsonValidation(req, request)

			if tc.expectedErr && validator.Valid() || !tc.expectedErr && !validator.Valid() {
				t.Errorf("unexpected error: %v", validator.Errors.MessageMap)
			} else if len(tc.expectedErrMap) != len(validator.Errors.MessageMap) {
				t.Errorf("unexpected error, want %v, but got %v", tc.expectedErrMap, validator.Errors.MessageMap)
			} else if !reflect.DeepEqual(tc.expectedErrMap, validator.Errors.MessageMap) {
				t.Errorf("expected error %s, but got %s", tc.expectedErrMap, validator.Errors.MessageMap)
			}
		})
	}
}
