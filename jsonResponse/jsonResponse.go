package jsonResponse

import (
	"encoding/json"
	"net/http"
)

func ResponseWithJson(w http.ResponseWriter, code int, payload interface{}) error {
	responseBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set(("Access-Control-Allow-Origin"), "*")
	w.WriteHeader(code)
	w.Write(responseBytes)
	return nil
}

func ResponseWithError(w http.ResponseWriter, code int, payload interface{}) error {
	return ResponseWithJson(w, code, payload)
}
