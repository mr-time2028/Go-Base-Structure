package json

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
)

// ReadJSON read request and extract request payload from it
func ReadJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	if r.Body == nil {
		return errors.New("request body is empty")
	}

	maxBytes := 1048576 // one megabyte
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(data)
	if err != nil {
		if err == io.EOF {
			return errors.New("request body is empty")
		}
		if syntaxErr, ok := err.(*json.SyntaxError); ok {
			return errors.New("request body contains invalid JSON syntax at position " + strconv.Itoa(int(syntaxErr.Offset)))
		}
		if unmarshalErr, ok := err.(*json.UnmarshalTypeError); ok {
			return errors.New("request body contains invalid data type at position " + strconv.Itoa(int(unmarshalErr.Offset)))
		}
		return errors.New("failed to decode request body: " + err.Error())
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
func ErrorJSON(w http.ResponseWriter, err error, status ...int) error {
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
		return err
	}

	return nil
}
