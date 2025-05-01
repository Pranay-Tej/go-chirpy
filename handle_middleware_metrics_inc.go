package main

import "net/http"

func (apiConfig *ApiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiConfig.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}
