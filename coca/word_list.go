// Library for converting the CSV wordlist file into a Go source file providing
// programmatic access to it without any runtime file dependencies.

package coca

import (
	"encoding/csv"
	"fmt"
	"github.com/sethpollen/dorkalonius"
	"io"
	"sort"
	"strconv"
	"strings"
)

// Fetches the COCA word list.
func GetWordList() (*dorkalonius.WordList, error) {
	reader := csv.NewReader(Get_coca_5000_csv("coca-5000.csv"))
	// Disable field count checking.
	reader.FieldsPerRecord = -1

	// Our raw data may contain 2 lines with the same word if that word can be
	// used as more than one part of speech. We just add the occurrence counts
	// of these lines together.
	var wordSet = make(map[string]*dorkalonius.Word)

	for i := 0; true; i++ {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if i < 2 {
			// Skip the first 2 lines; they are headers.
			continue
		}
		if len(record) != 5 {
			return nil, fmt.Errorf(
				"Wrong number of columns (", len(record), ") on line ", i)
		}

		occurrences, err := strconv.ParseInt(record[3], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Invalid occurrences on line ", i)
		}
		word := record[1]
		partOfSpeech := record[2]

		// Blacklist specific words we don't like from the data file.
		if word == "n't" {
			continue
		}

		if existing, found := wordSet[word]; found {
			existing.Occurrences += occurrences
			if strings.Index(existing.PartsOfSpeech, partOfSpeech) < 0 {
				existing.PartsOfSpeech += partOfSpeech
			}
		} else {
			wordSet[word] = &dorkalonius.Word{word, occurrences, partOfSpeech}
		}
	}

	// Convert the map to a WordList object.
	wordList := dorkalonius.NewWordList()
	for _, word := range wordSet {
		wordList.AddWord(*word)
	}

	sort.Sort(wordList)
	return wordList, nil
}
