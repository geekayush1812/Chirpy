package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/geekayush1812/chirpy/jsonResponse"
)

type userRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type updateUserRequest struct {
	Email       string `json:"email,omitempty"`
	Password    string `json:"password,omitempty"`
	IsChirpyRed bool   `json:"is_chirpy_red,omitempty"`
}

const (
	SOMETHING_WENT_WRONG        = "Something went wrong"
	REQUIRED_FIELD_ERROR        = "both email and password are required fields"
	ONE_OF_REQUIRED_FIELD_ERROR = "one of email and password are required"
	USER_NOT_AUTHORIZED         = "not authorized"
)

func handleUserRequestInputSanitization(w http.ResponseWriter, r *http.Request) (userRequest, error) {
	decoder := json.NewDecoder(r.Body)
	requestBody := userRequest{}
	err := decoder.Decode(&requestBody)

	if err != nil {
		handleApiError(w, http.StatusBadRequest, SOMETHING_WENT_WRONG)
		return userRequest{}, errors.New("")
	}

	if requestBody.Email == "" || requestBody.Password == "" {
		handleApiError(w, http.StatusUnprocessableEntity, REQUIRED_FIELD_ERROR)
		return userRequest{}, errors.New("")
	}

	return requestBody, nil
}

func createUser(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	requestBody, err := handleUserRequestInputSanitization(w, r)

	if err != nil {
		return
	}

	user, err := db.CreateUser(requestBody.Email, requestBody.Password)

	if err != nil {
		handleApiError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = jsonResponse.ResponseWithJson(w, http.StatusCreated, user)

	if err != nil {
		handleApiError(w, http.StatusBadRequest, SOMETHING_WENT_WRONG)
		return
	}
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	statusCode, errMessage, userJwtClaims := isUserAuthorized(w, r)

	if statusCode != http.StatusOK {
		handleApiError(w, statusCode, errMessage)
		return
	}

	decoder := json.NewDecoder(r.Body)
	requestBody := updateUserRequest{}

	err := decoder.Decode(&requestBody)

	if err != nil {
		handleApiError(w, http.StatusBadRequest, SOMETHING_WENT_WRONG)
		return
	}

	if requestBody.Email == "" && requestBody.Password == "" {
		handleApiError(w, http.StatusUnprocessableEntity, ONE_OF_REQUIRED_FIELD_ERROR)
		return
	}

	userId, err := strconv.ParseInt(userJwtClaims.Subject, 10, 32)

	if err != nil {
		handleApiError(w, http.StatusInternalServerError, SOMETHING_WENT_WRONG)
		return
	}

	user, err := db.UpdateUser(int(userId), requestBody.Email, requestBody.Password, false)

	if err != nil {
		handleApiError(w, http.StatusInternalServerError, SOMETHING_WENT_WRONG)
		return
	}

	jsonResponse.ResponseWithJson(w, http.StatusOK, user)
}
