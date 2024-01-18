package route

import (
	"dictionary-go/dictionary"
	"net/http"
	"sync"
)

func AddHandler(d *dictionary.Dictionary, wg *sync.WaitGroup, errors chan<- error) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		word := request.URL.Query().Get("word")
		definition := request.URL.Query().Get("definition")

		d.Add(word, definition, wg, errors)
	}
}
