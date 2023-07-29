package admin

import (
	"fmt"
	"net/http"

	"github.com/geekayush1812/chirpy/api"
	"github.com/go-chi/chi"
)

func GetAdminRouter(apiCfg *api.ApiConfig) http.Handler {

	adminRouter := chi.NewRouter()

	adminRouter.Get("/metrics", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", apiCfg.FileServerHits)))
	})

	return adminRouter
}
