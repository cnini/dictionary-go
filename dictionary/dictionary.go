package dictionary

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Dictionary struct {
	Filename string
}

type DictionaryEntry struct {
	Word       string
	Definition string
}

func NewDictionary(filename string) *Dictionary {
	return &Dictionary{
		Filename: filename,
	}
}

func (d Dictionary) Add(word string, definition string) error {
	dictionaryEntry := DictionaryEntry{
		Word:       word,
		Definition: definition,
	}

	dictionaryEntries, err := d.read()

	if err != nil {
		return err
	}

	dictionaryEntries = append(dictionaryEntries, dictionaryEntry)

	d.write(dictionaryEntries)

	return nil
}

func (d Dictionary) Get(searchTerm string) (string, string, error) {
	dictionaryEntries, err := d.read()

	if err != nil {
		return "", "", err
	}

	for _, dictionaryEntry := range dictionaryEntries {
		if dictionaryEntry.Word == searchTerm || dictionaryEntry.Definition == searchTerm {
			return dictionaryEntry.Word, dictionaryEntry.Definition, nil
		}
	}

	return "", "", err
}

// func (d Dictionary) Remove(termToRemove string) {
// 	// Get the word and not the definition,
// 	// even if termToRemove is definition
// 	word, _ := d.Get(termToRemove)

// 	delete(d, word)
// }

// func (d Dictionary) List() []string {
// 	var sortedDictionary []string

// 	// Create a list version of the dictionnary to easily sort it
// 	var listedDictionary []string
// 	for word := range d {
// 		listedDictionary = append(listedDictionary, word)
// 	}

// 	sort.Strings(listedDictionary)

// 	for _, word := range listedDictionary {
// 		definition := d[word]
// 		sortedDictionary = append(sortedDictionary, fmt.Sprintf("%s: %s", word, definition))
// 	}

// 	return sortedDictionary
// }

func (d Dictionary) read() ([]DictionaryEntry, error) {
	dictionaryFile, err := os.Open(d.Filename)

	if err != nil {
		return nil, err
	}

	defer dictionaryFile.Close()

	var dictionaryEntries []DictionaryEntry

	scanner := bufio.NewScanner(dictionaryFile)

	for scanner.Scan() {
		line := scanner.Text()

		// Split the line in two parts
		lineParts := strings.SplitN(line, ":", 2)

		if len(lineParts) == 2 {
			// Recreate a new DictionaryEntry
			dictionaryEntry := DictionaryEntry{
				Word:       strings.TrimSpace(lineParts[0]),
				Definition: strings.TrimSpace(lineParts[1]),
			}

			dictionaryEntries = append(dictionaryEntries, dictionaryEntry)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return dictionaryEntries, nil
}

func (d Dictionary) write(dictionaryEntries []DictionaryEntry) error {
	dictionaryFile, err := os.Create(d.Filename)

	if err != nil {
		return err
	}

	defer dictionaryFile.Close()

	for _, dictionaryEntry := range dictionaryEntries {
		line := fmt.Sprintf("%s: %s\n", dictionaryEntry.Word, dictionaryEntry.Definition)

		_, err := dictionaryFile.WriteString(line)

		if err != nil {
			return err
		}
	}

	return nil
}
