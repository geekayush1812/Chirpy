package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/geekayush1812/chirpy/jsonResponse"
	"github.com/go-chi/chi"
)

const (
	EMPTY_CHIRP_BODY_ERROR    = "Chirp is empty"
	TOO_LONG_CHIRP_BODY_ERROR = "Chirp is too long"
	CHIRP_ID_REQUIRED         = "Chirp id required"
	USER_FORBIDDEN            = "user not allowed"
)

type chirpRequest struct {
	Body string `json:"body"`
}

func getCleanedBody(body string) string {
	profaneWords := []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}
	splitBody := strings.Split(body, " ")
	for i, word := range splitBody {
		found := false
		for _, w := range profaneWords {
			if w == strings.ToLower(word) {
				found = true
				break
			}
		}
		if found {
			splitBody[i] = "****"
		}
	}

	return strings.Join(splitBody, " ")
}

func isChirpLengthValid(body string) bool {
	return len(body) <= 140 && len(body) > 0
}

func isChirpEmpty(body string) bool {
	return len(body) == 0
}

func createChirp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	statusCode, errMessage, userJwtClaims := isUserAuthorized(w, r)

	if statusCode != http.StatusOK {
		handleApiError(w, statusCode, errMessage)
		return
	}

	decoder := json.NewDecoder(r.Body)
	responseBody := chirpRequest{}

	err := decoder.Decode(&responseBody)

	if err != nil {
		handleApiError(w, http.StatusBadRequest, SOMETHING_WENT_WRONG)
		return
	}

	if isChirpEmpty(responseBody.Body) {
		handleApiError(w, http.StatusBadRequest, EMPTY_CHIRP_BODY_ERROR)
		return
	}

	if !isChirpLengthValid(responseBody.Body) {
		handleApiError(w, http.StatusBadRequest, TOO_LONG_CHIRP_BODY_ERROR)
		return
	}

	cleanedBody := getCleanedBody(responseBody.Body)

	userId, err := strconv.ParseInt(userJwtClaims.Subject, 10, 32)

	if err != nil {
		handleApiError(w, http.StatusInternalServerError, SOMETHING_WENT_WRONG)
		return
	}

	chirp, err := db.CreateChirp(int(userId), cleanedBody)

	if err != nil {
		handleApiError(w, http.StatusBadRequest, SOMETHING_WENT_WRONG)
		return
	}

	err = jsonResponse.ResponseWithJson(w, http.StatusCreated, chirp)

	if err != nil {
		handleApiError(w, http.StatusBadRequest, SOMETHING_WENT_WRONG)
		return
	}
}

func getChirpsById(w http.ResponseWriter, r *http.Request) {
	urlParams := chi.URLParam(r, "chirpId")

	if urlParams == "" {
		handleApiError(w, http.StatusBadRequest, CHIRP_ID_REQUIRED)
		return
	}

	chirpId, err := strconv.Atoi(urlParams)
	if err != nil {
		log.Fatal("could not convert chirpId to int")
		return
	}

	chirp, err := db.GetChirp(chirpId)

	if err != nil {
		jsonResponse.ResponseWithError(w, http.StatusNotFound, nil)
		return
	}

	err = jsonResponse.ResponseWithJson(w, http.StatusOK, chirp)

	if err != nil {
		handleApiError(w, http.StatusBadRequest, SOMETHING_WENT_WRONG)
		return
	}
}

func getChirps(w http.ResponseWriter, r *http.Request) {
	author_id := r.URL.Query().Get("author_id")
	sortType := r.URL.Query().Get("sort")

	chirps, err := db.GetChirps(author_id, sortType)

	if err != nil {
		handleApiError(w, http.StatusBadRequest, SOMETHING_WENT_WRONG)
		return
	}

	err = jsonResponse.ResponseWithJson(w, http.StatusOK, chirps)
	if err != nil {
		handleApiError(w, http.StatusBadRequest, SOMETHING_WENT_WRONG)
		return
	}
}

func deleteChirpById(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	statusCode, errMessage, userJwtClaims := isUserAuthorized(w, r)

	if statusCode != http.StatusOK {
		handleApiError(w, statusCode, errMessage)
		return
	}

	urlParams := chi.URLParam(r, "chirpId")

	if urlParams == "" {
		handleApiError(w, http.StatusBadRequest, CHIRP_ID_REQUIRED)
		return
	}

	chirpId, err := strconv.Atoi(urlParams)
	if err != nil {
		log.Fatal("could not convert chirpId to int")
		return
	}

	chirp, err := db.GetChirp(chirpId)

	if err != nil {
		handleApiError(w, http.StatusBadRequest, SOMETHING_WENT_WRONG)
		return
	}

	userId, err := strconv.ParseInt(userJwtClaims.Subject, 10, 32)

	if err != nil {
		handleApiError(w, http.StatusBadRequest, SOMETHING_WENT_WRONG)
		return
	}

	if chirp.AuthorId != int(userId) {
		handleApiError(w, http.StatusForbidden, USER_FORBIDDEN)
		return
	}

	err = db.DeleteChirp(chirpId)

	if err != nil {
		handleApiError(w, http.StatusForbidden, USER_FORBIDDEN)
		return
	}

	jsonResponse.ResponseWithJson(w, http.StatusOK, struct{}{})
}
