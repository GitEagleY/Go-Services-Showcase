package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type jsonResponce struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// readJSON reads a JSON request body from an HTTP request, decoding it into the provided data structure.
// It enforces a maximum body size of one megabyte (1 MB) and checks that the body contains only a single JSON value.
// If any errors occur during decoding or validation, it returns an error.
func (app *Config) readJSON(w http.ResponseWriter, r *http.Request, data interface{}) error {
	maxBytes := 1048576 // One megabyte (1 MB)

	// Enforce a maximum body size of 1 MB on the request body.
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// Create a JSON decoder for the request body.
	dec := json.NewDecoder(r.Body)

	// Decode the JSON data into the provided data structure.
	err := dec.Decode(data)
	if err != nil {
		return err
	}

	// Check that the body contains only a single JSON value (EOF).
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must have only a single JSON value")
	}

	return nil
}

// writeJSON writes a JSON response to an HTTP response writer with the specified HTTP status code.
// It serializes the provided data structure to JSON and sets the "Content-Type" header to "application/json".
// Additional headers can be provided as a variadic argument. It returns an error if serialization or writing fails.
func (app *Config) writeJSON(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	// Serialize the data structure to JSON.
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Set the HTTP response headers.
	if len(headers) > 0 {
		for key, value := range headers[0] {
			w.Header()[key] = value
		}
	}

	// Set the "Content-Type" header to indicate that the response is JSON.
	w.Header().Set("Content-Type", "application/json")

	// Set the HTTP status code.
	w.WriteHeader(status)

	// Write the JSON response to the response writer.
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

// errorJSON is a convenience function for writing JSON error responses.
// It takes an error message, an optional HTTP status code (defaulting to 400 Bad Request), and constructs a JSON response payload.
// It then delegates to the writeJSON function to write the error response.
func (app *Config) errorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	// If an optional status code is provided, use it.
	if len(status) > 0 {
		statusCode = status[0]
	}

	// Create a JSON response payload with an error message.
	payload := jsonResponce{
		Error:   true,
		Message: err.Error(),
	}

	// Delegate to writeJSON to write the error response.
	return app.writeJSON(w, statusCode, payload)
}
