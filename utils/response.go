package utils

import (
	"encoding/json"
	"net/http"
)

// Return success response
func Success(w http.ResponseWriter, status int, resp map[string]interface{}, data interface{}, message string) {
	resp["success"] = true
	resp["data"] = data
	resp["message"] = message
	Respond(w, status, resp)
}

// Return fail response
func Fail(w http.ResponseWriter, status int, resp map[string]interface{}, message string) {
	resp["error"] = message
	Respond(w, status, resp)
}

// Return json response
func Respond(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Add("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
