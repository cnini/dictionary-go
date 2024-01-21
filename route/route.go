package route

import (
	"dictionary-go/dictionary"
	"dictionary-go/middleware"
	"fmt"
	"net/http"
	"sync"
)

func AddHandler(d *dictionary.Dictionary, wg *sync.WaitGroup, errors chan<- error) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if middleware.IsAuthorized(request) {
			word := request.URL.Query().Get("word")
			definition := request.URL.Query().Get("definition")

			middleware.NewLog(
				func(w http.ResponseWriter, r *http.Request) {
					d.Add(word, definition, wg, errors)
				},
				writer,
				request,
				errors,
				"AUTHORIZED",
			)

			fmt.Println("\n-- French-English dictionary ------")
			sortedDictionary := d.List(errors)

			for _, entry := range sortedDictionary {
				fmt.Println(entry)
			}
		} else {
			middleware.NewLog(
				func(w http.ResponseWriter, r *http.Request) {
					http.Error(writer, "Permission denied", http.StatusForbidden)
				},
				writer,
				request,
				errors,
				"UNAUTHORIZED",
			)
		}
	}
}

func GetHandler(d *dictionary.Dictionary, errors chan<- error) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		searchTerm := request.URL.Query().Get("word")

		middleware.NewLog(
			func(w http.ResponseWriter, r *http.Request) {
				word, definition := d.Get(searchTerm, errors)

				if word != "" {
					fmt.Printf("\n-- \"%s\": %s. ------\n", word, definition)
				} else {
					fmt.Printf("\n-- \"%s\" does not exists. ------\n", searchTerm)
				}
			},
			writer,
			request,
			errors,
			"AUTHORIZED",
		)
	}
}

func RemoveHandler(d *dictionary.Dictionary, wg *sync.WaitGroup, errors chan<- error) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		if middleware.IsAuthorized(request) {
			word := request.URL.Query().Get("word")

			middleware.NewLog(
				func(w http.ResponseWriter, r *http.Request) {
					d.Remove(word, wg, errors)
				},
				writer,
				request,
				errors,
				"AUTHORIZED",
			)

			fmt.Println("\n-- French-English dictionary ------")
			sortedDictionary := d.List(errors)

			for _, entry := range sortedDictionary {
				fmt.Println(entry)
			}
		} else {
			middleware.NewLog(
				func(w http.ResponseWriter, r *http.Request) {
					http.Error(writer, "Permission denied", http.StatusForbidden)
				},
				writer,
				request,
				errors,
				"UNAUTHORIZED",
			)
		}
	}
}
