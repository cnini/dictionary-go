package main

import (
	"dictionary-go/dictionary"
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	errors := make(chan error)

	defer close(errors)

	dictionary := dictionary.NewDictionary("dictionary.txt", &wg, errors)

	dictionary.Add("Apprendre", "To learn", &wg, errors)
	dictionary.Add("Enseigner", "To teach", &wg, errors)
	dictionary.Add("Livre", "Book", &wg, errors)
	dictionary.Add("Crayon", "Pencil", &wg, errors)
	dictionary.Add("RAS", "Nothing to report", &wg, errors)
	dictionary.Add("MDR", "LOL", &wg, errors)
	dictionary.Add("Supprimer", "To remove", &wg, errors)

	// dictionaryEntries := []struct {
	// 	Word       string
	// 	Definition string
	// }{
	// 	{"Apprendre", "To learn"},
	// 	{"Enseigner", "To teach"},
	// 	{"Livre", "Book"},
	// 	{"Crayon", "Pencil"},
	// 	{"RAS", "Nothing to report"},
	// 	{"MDR", "LOL"},
	// 	{"Supprimer", "To remove"},
	// }

	// for _, dictionaryEntry := range dictionaryEntries {
	// 	dictionary.Add(dictionaryEntry.Word, dictionaryEntry.Definition, &wg, errors)
	// }

	fmt.Println("-- French-English dictionary (after Add calls) ------")
	sortedDictionary := dictionary.List(&wg, errors)

	for _, entry := range sortedDictionary {
		fmt.Println(entry)
	}

	searchTerm := "RAS"
	word, definition := dictionary.Get(searchTerm, &wg, errors)

	if word != "" {
		fmt.Printf("\n-- \"%s\"'s definition : %s. ------", word, definition)
	} else {
		fmt.Printf("\n-- \"%s\" does not exists. ------", searchTerm)
	}

	termToRemove := "To remove"
	fmt.Printf("\n\n-- Removing \"%s\" line. ------", termToRemove)
	dictionary.Remove(termToRemove, &wg, errors)

	fmt.Println("\n\n-- French-English dictionary (after Remove call) ------")
	sortedDictionary = dictionary.List(&wg, errors)

	for _, entry := range sortedDictionary {
		fmt.Println(entry)
	}
}
