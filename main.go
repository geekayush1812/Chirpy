package main

import (
	"log"
	"net/http"

	"github.com/geekayush1812/chirpy/admin"
	"github.com/geekayush1812/chirpy/api"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	fileRootPath := "."
	port := ":8080"

	mainRouter := chi.NewRouter()
	apiRouter, apiCfg := api.GetApiRouter()
	adminRouter := admin.GetAdminRouter(apiCfg)

	mainRouter.Mount("/", apiCfg.MiddlewareMetricsInc(http.FileServer(http.Dir(fileRootPath))))
	mainRouter.Mount("/api", apiRouter)
	mainRouter.Mount("/admin", adminRouter)

	server := http.Server{
		Handler: middlewareCors(mainRouter),
		Addr:    port,
	}

	log.Printf("Serving on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}
