package main

import (
	"github.com/ylapshin/urlshortener/internal/app"
	"log"
	"net/http"
)

func main() {
	app := app.New()
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, app.RootHandler)
	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		log.Panicln(err)
	}
}
