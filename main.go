package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/domain", DomaineHandler).Methods("GET")

	http.ListenAndServe(":8088", r)
}
