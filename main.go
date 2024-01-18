package main

import (
	"dictionary-go/dictionary"
	"fmt"
)

func main() {
	dictionary := dictionary.NewDictionary("dictionary.txt")

	dictionary.Add("Apprendre", "To learn")
	dictionary.Add("Enseigner", "To teach")
	dictionary.Add("Livre", "Book")
	dictionary.Add("Crayon", "Pencil")
	dictionary.Add("RAS", "Nothing to report")
	dictionary.Add("MDR", "LOL")
	dictionary.Add("Supprimer", "To remove")

	// fmt.Println("-- French-English dictionary (after Add calls) ------")
	// sortedDictionary := dictionary.List()
	// for _, entry := range sortedDictionary {
	// 	fmt.Println(entry)
	// }

	searchTerm := "RAS"
	word, definition, getErr := dictionary.Get(searchTerm)

	if getErr != nil {
		fmt.Printf("\n-- Erreur : %s ------", getErr.Error())
	} else {
		if word != "" {
			fmt.Printf("\n-- \"%s\"'s definition : %s. ------", word, definition)
		} else {
			fmt.Printf("\n-- \"%s\" does not exists. ------", searchTerm)
		}
	}

	termToRemove := "To remove"
	fmt.Printf("\n\n-- Removing \"%s\" line. ------", termToRemove)
	dictionary.Remove(termToRemove)

	// fmt.Println("\n\n-- French-English dictionary (after Remove call) ------")
	// sortedDictionary = dictionary.List()
	// for _, entry := range sortedDictionary {
	// 	fmt.Println(entry)
	// }
}
