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

func NewDictionary(filename string, wg *sync.WaitGroup, errors chan<- error) *Dictionary {
	wg.Add(1)

	dictionary := &Dictionary{
		Filename: filename,
	}

	dictionary.clearFile(wg, errors)

	return dictionary
}

func (d Dictionary) Add(word string, definition string, wg *sync.WaitGroup, errors chan<- error) {
	wg.Add(1)

	dictionaryEntry := DictionaryEntry{
		Word:       word,
		Definition: definition,
	}

	dictionaryEntries := d.read(wg, errors)

	wg.Add(1)

	dictionaryEntries = append(dictionaryEntries, dictionaryEntry)

	d.write(dictionaryEntries, wg, errors)
}

func (d Dictionary) Get(searchTerm string, wg *sync.WaitGroup, errors chan<- error) (string, string) {
	wg.Add(1)

	dictionaryEntries := d.read(wg, errors)

	for _, dictionaryEntry := range dictionaryEntries {
		if dictionaryEntry.Word == searchTerm || dictionaryEntry.Definition == searchTerm {
			return dictionaryEntry.Word, dictionaryEntry.Definition
		}
	}

	return "", ""
}

func (d Dictionary) Remove(termToRemove string, wg *sync.WaitGroup, errors chan<- error) {
	wg.Add(1)

	dictionaryEntries := d.read(wg, errors)

	var updatedDictionaryEntries []DictionaryEntry
	for _, dictionaryEntry := range dictionaryEntries {
		if dictionaryEntry.Word != termToRemove && dictionaryEntry.Definition != termToRemove {
			updatedDictionaryEntries = append(updatedDictionaryEntries, dictionaryEntry)
		}
	}

	wg.Add(1)

	d.write(updatedDictionaryEntries, wg, errors)
}

func (d Dictionary) List(wg *sync.WaitGroup, errors chan<- error) []string {
	wg.Add(1)

	var sortedDictionary []string

	// Create a list version of the dictionnary to easily sort it
	var listedDictionary []string

	dictionaryEntries := d.read(wg, errors)

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

func (d *Dictionary) clearFile(wg *sync.WaitGroup, errors chan<- error) {
	defer wg.Done()

	file, err := os.Create(d.Filename)
	if err != nil {
		errors <- err
	}

	defer file.Close()
}

func (d Dictionary) read(wg *sync.WaitGroup, errors chan<- error) []DictionaryEntry {
	defer wg.Done()

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

func (d Dictionary) write(dictionaryEntries []DictionaryEntry, wg *sync.WaitGroup, errors chan<- error) {
	defer wg.Done()

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
