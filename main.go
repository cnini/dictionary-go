package main

import (
	"dictionary-go/server"
)

func main() {
	server.Start()
	server.ListenAndServe()

	// fmt.Println("-- French-English dictionary (after Add calls) ------")
	// sortedDictionary := dictionary.List(errors)

	// for _, entry := range sortedDictionary {
	// 	fmt.Println(entry)
	// }

	// searchTerm := "RAS"
	// word, definition := dictionary.Get(searchTerm, errors)

	// if word != "" {
	// 	fmt.Printf("\n-- \"%s\"'s definition : %s. ------", word, definition)
	// } else {
	// 	fmt.Printf("\n-- \"%s\" does not exists. ------", searchTerm)
	// }

	// fmt.Println("\n\n-- French-English dictionary (after Remove call) ------")
	// sortedDictionary = dictionary.List(errors)

	// for _, entry := range sortedDictionary {
	// 	fmt.Println(entry)
	// }
}
