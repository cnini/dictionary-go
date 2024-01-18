package server

import (
	"dictionary-go/dictionary"
	"dictionary-go/route"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

func Start() {
	var wg sync.WaitGroup
	errors := make(chan error)
	defer close(errors)

	dictionary := dictionary.NewDictionary("dictionary.txt", errors)

	router := mux.NewRouter()
	http.Handle("/", router)

	router.HandleFunc("/add", route.AddHandler(dictionary, &wg, errors)).Methods("POST")
}

func ListenAndServe() {
	http.ListenAndServe(":8080", nil)
}
