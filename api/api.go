package api

import (
	"net/http"

	"github.com/geekayush1812/chirpy/database"
	"github.com/go-chi/chi"
)

type ApiConfig struct {
	FileServerHits int
}

type apiResponseError struct {
	Error string `json:"error"`
}

var dbPath = "database.json"
var db = database.NewDB(dbPath)

func GetApiRouter() (http.Handler, *ApiConfig) {

	apiRouter := chi.NewRouter()
	apiCfg := &ApiConfig{
		FileServerHits: 0,
	}

	apiRouter.Get("/healthz", healthz)

	apiRouter.Get("/metrics", apiCfg.middlewareMetrics())

	apiRouter.Group(func(r chi.Router) {
		r.Post("/chirps", createChirp)
		r.Get("/chirps", getChirps)
		r.Get("/chirps/{chirpId}", getChirpsById)
		r.Delete("/chirps/{chirpId}", deleteChirpById)
	})

	apiRouter.Post("/users", createUser)
	apiRouter.Put("/users", updateUser)

	apiRouter.Post("/login", login)

	apiRouter.Post("/refresh", refreshToken)
	apiRouter.Post("/revoke", revokeToken)

	apiRouter.Post("/polka/webhooks", handlePolkaWebHooks)

	return apiRouter, apiCfg
}
