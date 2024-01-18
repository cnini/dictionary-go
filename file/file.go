package file

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/cnini/dictionary-go/dictionary"
)

func Read(filename string) ([]dictionary.DictionaryEntry, error) {
	dictionaryFile, err := os.Open(filename)

	if err != nil {
		return nil, err
	}

	defer dictionaryFile.Close()

	var dictionaryEntries []dictionary.DictionaryEntry

	scanner := bufio.NewScanner(dictionaryFile)

	for scanner.Scan() {
		line := scanner.Text()

		// Split the line in two parts
		lineParts := strings.SplitN(line, ":", 2)

		if len(lineParts) == 2 {
			// Recreate a new DictionaryEntry
			dictionaryEntry := dictionary.DictionaryEntry{
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

func Write(filename string, dictionaryEntries []dictionary.DictionaryEntry) error {
	dictionaryFile, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)

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
