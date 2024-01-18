package dictionary

import (
	"bufio"
	"fmt"
	"os"
	"sort"
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
	dictionary := &Dictionary{
		Filename: filename,
	}

	dictionary.clearFile()

	return dictionary
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

func (d Dictionary) Remove(termToRemove string) error {
	dictionaryEntries, err := d.read()
	if err != nil {
		return err
	}

	var updatedDictionaryEntries []DictionaryEntry
	for _, dictionaryEntry := range dictionaryEntries {
		if dictionaryEntry.Word != termToRemove && dictionaryEntry.Definition != termToRemove {
			updatedDictionaryEntries = append(updatedDictionaryEntries, dictionaryEntry)
		}
	}

	d.write(updatedDictionaryEntries)

	return nil
}

func (d Dictionary) List() ([]string, error) {
	var sortedDictionary []string

	// Create a list version of the dictionnary to easily sort it
	var listedDictionary []string
	dictionaryEntries, err := d.read()
	if err != nil {
		return nil, err
	}

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

	return sortedDictionary, nil
}

func (d *Dictionary) clearFile() error {
	file, err := os.Create(d.Filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return nil
}

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
