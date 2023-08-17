package helpers

import (
	"encoding/json"
	"errors"
	"go-base-structure/cmd/config"
	"io"
	"net/http"
	"os"
	"strconv"
)

// ReadJSON read request and extract request payload from it
func ReadJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	maxBytes := 1048576 // one megabyte
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must have only a single json value")
	}

	return nil
}

// WriteJSON write data to output
func WriteJSON(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

// ErrorJSON write an error message to output
func ErrorJSON(w http.ResponseWriter, err error, status ...int) {
	type jsonResponse struct {
		Error   bool   `json:"error"`
		Message string `json:"message"`
	}

	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload jsonResponse
	payload.Error = true
	payload.Message = err.Error()

	err = WriteJSON(w, statusCode, payload)
	if err != nil {
		config.AppConfig.ErrorLog.Println(err)
	}
}

// GetEnvOrDefaultString read string data from env file
func GetEnvOrDefaultString(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetEnvOrDefaultBool read bool data from env file
func GetEnvOrDefaultBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	boolVal, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return boolVal
}

// GetEnvOrDefaultInt read int data from env file
func GetEnvOrDefaultInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intVal, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return intVal
}

// GetEnvOrDefaultFloat32 read float32 type from env file
func GetEnvOrDefaultFloat32(key string, defaultValue float32) float32 {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	float32Val, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return defaultValue
	}
	return float32(float32Val)
}

// GetEnvOrDefaultFloat64 read float64 file from env file
func GetEnvOrDefaultFloat64(key string, defaultValue float64) float64 {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	float64Val, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return defaultValue
	}
	return float64Val
}
