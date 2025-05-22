package utils

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func DecodeJson(r *http.Request, obj interface{}) error {
	if r.Body == nil {
		return errors.New("request body is empty")
	}
	defer r.Body.Close()

	limitedReader := io.LimitReader(r.Body, 1048576*5) // Limit the size of the request body at 5 MB

	decoder := json.NewDecoder(limitedReader)
	decoder.DisallowUnknownFields() // Prevent silently ignoring unexpected fields

	if err := decoder.Decode(obj); err != nil {
		return err
	}

	if decoder.More() {
		return errors.New("request body must contain only a single JSON object")
	}

	return nil
}

func DncodeJson(w http.ResponseWriter, statusCode int, obj interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(obj); err != nil {
		http.Error(w, `{"error":"failed to encode response"}`, http.StatusInternalServerError)
		return err
	}

	return nil
}
