package api

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/geekayush1812/chirpy/jsonResponse"
	"github.com/golang-jwt/jwt/v5"
)

type UserToken struct {
	Token string `json:"token"`
}

func refreshToken(w http.ResponseWriter, r *http.Request) {
	authorizationHeader := r.Header.Get("Authorization")

	if authorizationHeader == "" {
		handleApiError(w, http.StatusUnauthorized, USER_NOT_AUTHORIZED)
		return
	}

	userBearerToken := strings.Replace(authorizationHeader, "Bearer ", "", 1)

	jwtSecret := os.Getenv("JWT_SECRET")

	userJwtClaims := &jwt.RegisteredClaims{}

	claims, err := jwt.ParseWithClaims(userBearerToken, userJwtClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil || !claims.Valid {
		handleApiError(w, http.StatusUnauthorized, USER_NOT_AUTHORIZED)
		return
	}

	userJwtIssuer := userJwtClaims.Issuer

	if userJwtIssuer != "chirpy-refresh" {
		handleApiError(w, http.StatusUnauthorized, USER_NOT_AUTHORIZED)
		return
	}

	userId, err := strconv.ParseInt(userJwtClaims.Subject, 10, 32)

	if err != nil {
		handleApiError(w, http.StatusInternalServerError, SOMETHING_WENT_WRONG)
		return
	}

	err = db.IsUserRefreshTokenRevoked(int(userId), userBearerToken)

	if err != nil {
		handleApiError(w, http.StatusUnauthorized, err.Error())
		return
	}

	unsignedJwtAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy-access",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		Subject:   strconv.FormatInt(userId, 10),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
	})

	signedJwtAccessToken, errAccessToken := unsignedJwtAccessToken.SignedString([]byte(jwtSecret))

	if errAccessToken != nil {
		jsonResponse.ResponseWithError(w, http.StatusInternalServerError, apiResponseError{
			Error: "something went wrong",
		})
		return
	}

	jsonResponse.ResponseWithJson(w, http.StatusOK, UserToken{
		Token: signedJwtAccessToken,
	})
}

func revokeToken(w http.ResponseWriter, r *http.Request) {
	authorizationHeader := r.Header.Get("Authorization")

	if authorizationHeader == "" {
		handleApiError(w, http.StatusUnauthorized, USER_NOT_AUTHORIZED)
		return
	}

	userBearerToken := strings.Replace(authorizationHeader, "Bearer ", "", 1)

	jwtSecret := os.Getenv("JWT_SECRET")

	userJwtClaims := &jwt.RegisteredClaims{}

	claims, err := jwt.ParseWithClaims(userBearerToken, userJwtClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil || !claims.Valid {
		handleApiError(w, http.StatusUnauthorized, USER_NOT_AUTHORIZED)
		return
	}

	userJwtIssuer := userJwtClaims.Issuer

	if userJwtIssuer != "chirpy-refresh" {
		handleApiError(w, http.StatusUnauthorized, USER_NOT_AUTHORIZED)
		return
	}

	userId, err := strconv.ParseInt(userJwtClaims.Subject, 10, 32)

	if err != nil {
		handleApiError(w, http.StatusInternalServerError, SOMETHING_WENT_WRONG)
		return
	}

	err = db.RevokeUserRefreshToken(int(userId), userBearerToken)

	if err != nil {
		handleApiError(w, http.StatusInternalServerError, SOMETHING_WENT_WRONG)
		return
	}

	jsonResponse.ResponseWithJson(w, http.StatusOK, UserToken{
		Token: "",
	})
}
