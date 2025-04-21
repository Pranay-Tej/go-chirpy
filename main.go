package main

import "net/http"

func main() {
	mux := http.NewServeMux()
	httpServer := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	httpServer.ListenAndServe()
}
