package common

import (
	"encoding/json"
	"net/http"
)

// Helper function for adding JSON data into a HTTP response.
func WriteJSON(
  w      http.ResponseWriter,
  status int,
  data   any,
) error {
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(status)
  return json.NewEncoder(w).Encode(data)
}

// Helper function for reading JSON data from a HTTP response.
func ReadJSON(
  r    *http.Request,
  data interface{},
) error {
  return json.NewDecoder(r.Body).Decode(data)
}

// Helper function for writing an error into a HTTP response.
func WriteError(
  w       http.ResponseWriter,
  status  int,
  message string,
) error {
  return WriteJSON(w, status, map[string]string{
    "error": message,
  })
}
