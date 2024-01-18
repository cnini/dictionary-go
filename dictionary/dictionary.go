package dictionary

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
)

type Dictionary struct {
	Filename string
}

type DictionaryEntry struct {
	Word       string
	Definition string
}

func NewDictionary(filename string, errors chan<- error) *Dictionary {
	dictionary := &Dictionary{
		Filename: filename,
	}

	dictionary.clearFile(errors)

	return dictionary
}

func (d Dictionary) Add(word string, definition string, wg *sync.WaitGroup, errors chan<- error) {
	wg.Add(1)
	defer wg.Done()

	dictionaryEntry := DictionaryEntry{
		Word:       word,
		Definition: definition,
	}

	dictionaryEntries := d.read(errors)

	dictionaryEntries = append(dictionaryEntries, dictionaryEntry)

	d.write(dictionaryEntries, errors)
}

func (d Dictionary) Get(searchTerm string, errors chan<- error) (string, string) {
	dictionaryEntries := d.read(errors)

	for _, dictionaryEntry := range dictionaryEntries {
		if dictionaryEntry.Word == searchTerm || dictionaryEntry.Definition == searchTerm {
			return dictionaryEntry.Word, dictionaryEntry.Definition
		}
	}

	return "", ""
}

func (d Dictionary) Remove(termToRemove string, wg *sync.WaitGroup, errors chan<- error) {
	wg.Add(1)
	defer wg.Done()

	dictionaryEntries := d.read(errors)

	var updatedDictionaryEntries []DictionaryEntry
	for _, dictionaryEntry := range dictionaryEntries {
		if dictionaryEntry.Word != termToRemove && dictionaryEntry.Definition != termToRemove {
			updatedDictionaryEntries = append(updatedDictionaryEntries, dictionaryEntry)
		}
	}

	d.write(updatedDictionaryEntries, errors)
}

func (d Dictionary) List(errors chan<- error) []string {
	var sortedDictionary []string

	// Create a list version of the dictionnary to easily sort it
	var listedDictionary []string

	dictionaryEntries := d.read(errors)

	for _, dictionaryEntry := range dictionaryEntries {
		listedDictionary = append(listedDictionary, dictionaryEntry.Word)
	}

	sort.Strings(listedDictionary)

	for _, word := range listedDictionary {
		for _, dictionaryEntry := range dictionaryEntries {
			if dictionaryEntry.Word == word {
				sortedDictionary = append(
					sortedDictionary,
					fmt.Sprintf("%s: %s", dictionaryEntry.Word, dictionaryEntry.Definition),
				)
			}
		}
	}

	return sortedDictionary
}

func (d *Dictionary) clearFile(errors chan<- error) {
	file, err := os.Create(d.Filename)
	if err != nil {
		errors <- err
	}

	defer file.Close()
}

func (d Dictionary) read(errors chan<- error) []DictionaryEntry {
	dictionaryFile, err := os.Open(d.Filename)
	if err != nil {
		errors <- err
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
		errors <- err
	}

	return dictionaryEntries
}

func (d Dictionary) write(dictionaryEntries []DictionaryEntry, errors chan<- error) {
	dictionaryFile, err := os.Create(d.Filename)
	if err != nil {
		errors <- err
	}

	defer dictionaryFile.Close()

	for _, dictionaryEntry := range dictionaryEntries {
		line := fmt.Sprintf("%s: %s\n", dictionaryEntry.Word, dictionaryEntry.Definition)

		_, err := dictionaryFile.WriteString(line)
		if err != nil {
			errors <- err
		}
	}
}
