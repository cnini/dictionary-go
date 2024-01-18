package dictionary

import (
	"fmt"
	"sort"
)

type Dictionary map[string]string

func (d Dictionary) Add(word string, definition string) {
	d[word] = definition
}

func (d Dictionary) Get(searchTerm string) (string, string) {
	// If searchTerm is not a definition
	if definition, found := d[searchTerm]; found {
		return searchTerm, definition
	}

	// If searchTerm is a definition
	for word, definition := range d {
		if definition == searchTerm {
			return word, definition
		}
	}

	return "", ""
}

func (d Dictionary) Remove(termToRemove string) {
	// Get the word and not the definition,
	// even if termToRemove is definition
	word, _ := d.Get(termToRemove)

	delete(d, word)
}

func (d Dictionary) List() []string {
	var sortedDictionary []string

	// Create a list version of the dictionnary to easily sort it
	var listedDictionary []string
	for word := range d {
		listedDictionary = append(listedDictionary, word)
	}

	sort.Strings(listedDictionary)

	for _, word := range listedDictionary {
		definition := d[word]
		sortedDictionary = append(sortedDictionary, fmt.Sprintf("%s: %s", word, definition))
	}

	return sortedDictionary
}
