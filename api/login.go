package api

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/geekayush1812/chirpy/jsonResponse"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserResponse struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	IsChirpyRed bool `json:"is_chirpy_red"`
}


func login(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()
	
	requestBody, err := handleUserRequestInputSanitization(w, r)

	if err != nil {
		return
	}

	user, err := db.GetUser(requestBody.Email)

	if err != nil {
		jsonResponse.ResponseWithError(w, http.StatusBadRequest, apiResponseError{
			Error: err.Error(),
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestBody.Password))

	if err != nil {
		jsonResponse.ResponseWithError(w, http.StatusUnauthorized, apiResponseError{
			Error: "incorrect credentials",
		})
		return
	}

	unsignedJwtAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer: "chirpy-access",
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		Subject: strconv.FormatInt(int64(user.Id), 10),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
	})

	unsignedJwtRefreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer: "chirpy-refresh",
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		Subject: strconv.FormatInt(int64(user.Id), 10),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(60 * 24 * time.Hour)),
	})

	jwtSecret := os.Getenv("JWT_SECRET")

	signedJwtAccessToken, errAccessToken := unsignedJwtAccessToken.SignedString([]byte(jwtSecret))
	signedJwtRefreshToken, errRefreshToken := unsignedJwtRefreshToken.SignedString([]byte(jwtSecret))

	if errAccessToken != nil || errRefreshToken != nil {
		jsonResponse.ResponseWithError(w, http.StatusInternalServerError, apiResponseError{
			Error: "something went wrong",
		})
		return
	}

	jsonResponse.ResponseWithJson(w, http.StatusOK, UserResponse{
		Id: user.Id,
		Email: user.Email,
		Token: signedJwtAccessToken,
		RefreshToken: signedJwtRefreshToken,
		IsChirpyRed: user.IsChirpyRed,
	})
}