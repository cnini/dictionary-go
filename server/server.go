package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Start() {
	router := mux.NewRouter()

	http.Handle("/", router)
}

func ListenAndServe() {
	http.ListenAndServe(":8080", nil)
}
