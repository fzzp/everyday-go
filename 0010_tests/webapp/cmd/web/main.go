package main

import (
	"log"
	"net/http"
)

const addr = ":8631"

type application struct{}

func main() {
	app := application{}

	mux := app.routes()

	log.Println("Starting server on port ", addr)
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		log.Fatal(err)
	}
}
