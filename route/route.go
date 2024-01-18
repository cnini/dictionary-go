package route

import (
	"dictionary-go/dictionary"
	"fmt"
	"net/http"
	"sync"
)

func AddHandler(d *dictionary.Dictionary, wg *sync.WaitGroup, errors chan<- error) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		word := request.URL.Query().Get("word")
		definition := request.URL.Query().Get("definition")

		d.Add(word, definition, wg, errors)

		fmt.Println("\n-- French-English dictionary ------")
		sortedDictionary := d.List(errors)

		for _, entry := range sortedDictionary {
			fmt.Println(entry)
		}
	}
}

func GetHandler(d *dictionary.Dictionary, errors chan<- error) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		searchTerm := request.URL.Query().Get("word")

		word, definition := d.Get(searchTerm, errors)

		if word != "" {
			fmt.Printf("\n-- \"%s\": %s. ------\n", word, definition)
		} else {
			fmt.Printf("\n-- \"%s\" does not exists. ------\n", searchTerm)
		}
	}
}

func RemoveHandler(d *dictionary.Dictionary, wg *sync.WaitGroup, errors chan<- error) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		word := request.URL.Query().Get("word")

		d.Remove(word, wg, errors)

		fmt.Println("\n-- French-English dictionary ------")
		sortedDictionary := d.List(errors)

		for _, entry := range sortedDictionary {
			fmt.Println(entry)
		}
	}
}
