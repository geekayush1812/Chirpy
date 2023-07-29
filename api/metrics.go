package api

import (
	"fmt"
	"net/http"
)

func (c *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.FileServerHits++
		next.ServeHTTP(w, r)
	})
}

func (c *ApiConfig) middlewareMetrics() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(fmt.Sprintf("Hits: %d", c.FileServerHits)))
	})
}
