package api

import (
	"net/http"

	"github.com/geekayush1812/chirpy/jsonResponse"
)


func handleApiError(w http.ResponseWriter, code int, err string) error {
	errResponse := apiResponseError{
		Error: err,
	}

	return jsonResponse.ResponseWithError(w, code, errResponse)
}