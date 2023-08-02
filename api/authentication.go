package api

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func isUserAuthorized(w http.ResponseWriter, r *http.Request) (int, string, jwt.RegisteredClaims) {
	authorizationHeader := r.Header.Get("Authorization")

	if authorizationHeader == "" {
		return http.StatusUnauthorized, USER_NOT_AUTHORIZED, jwt.RegisteredClaims{}
	}

	userBearerToken := strings.Replace(authorizationHeader, "Bearer ", "", 1)

	jwtSecret := os.Getenv("JWT_SECRET")

	userJwtClaims := &jwt.RegisteredClaims{}

	claims, err := jwt.ParseWithClaims(userBearerToken, userJwtClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil || !claims.Valid {
		return http.StatusUnauthorized, USER_NOT_AUTHORIZED, jwt.RegisteredClaims{}
	}

	userJwtExpiresAt := userJwtClaims.ExpiresAt
	userJwtIssuer := userJwtClaims.Issuer

	if userJwtIssuer != "chirpy-access" {
		return http.StatusUnauthorized, USER_NOT_AUTHORIZED, jwt.RegisteredClaims{}
	}

	isExpired := userJwtExpiresAt.Before(time.Now())

	if isExpired {
		return http.StatusUnauthorized, USER_NOT_AUTHORIZED, jwt.RegisteredClaims{}
	}

	return http.StatusOK, "", *userJwtClaims
}
