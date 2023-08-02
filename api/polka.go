package api

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/geekayush1812/chirpy/jsonResponse"
)

type PolkaWebhooksRequest struct {
	Event string `json:"event"`
	Data  struct {
		UserId int `json:"user_id"`
	} `json:"data"`
}

func isApiKeyValid(authorizationHeader string) bool {
	apiKey := strings.Replace(authorizationHeader, "ApiKey ", "", 1)

	polkaApiKey := os.Getenv("POLKA_API_KEY")

	return apiKey == polkaApiKey
}

func handlePolkaWebHooks(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	authorizationHeader := r.Header.Get("Authorization")

	if authorizationHeader == "" || !isApiKeyValid(authorizationHeader) {
		handleApiError(w, http.StatusUnauthorized, USER_NOT_AUTHORIZED)
		return
	}

	decoder := json.NewDecoder(r.Body)
	polkaWebhookRequest := PolkaWebhooksRequest{}

	err := decoder.Decode(&polkaWebhookRequest)

	if err != nil {
		handleApiError(w, http.StatusBadRequest, SOMETHING_WENT_WRONG)
		return
	}

	if polkaWebhookRequest.Event != "user.upgraded" {
		jsonResponse.ResponseWithJson(w, http.StatusOK, nil)
		return
	}

	// update user's is_chirpy_red flag
	err, userExists := db.IsUserExists(polkaWebhookRequest.Data.UserId)

	if err != nil {
		handleApiError(w, http.StatusBadRequest, SOMETHING_WENT_WRONG)
		return
	}

	if !userExists {
		handleApiError(w, http.StatusNotFound, "User not found")
		return
	}

	_, err = db.UpdateUser(polkaWebhookRequest.Data.UserId, "", "", true)

	if err != nil {
		handleApiError(w, http.StatusBadRequest, SOMETHING_WENT_WRONG)
		return
	}

	jsonResponse.ResponseWithJson(w, http.StatusOK, nil)
}
