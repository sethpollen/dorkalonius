// Library for converting the CSV wordlist file into a Go source file providing
// programmatic access to it without any runtime file dependencies.

package dorkalonius

import (
	"encoding/csv"
	"io"
	"log"
	"strconv"
	"strings"
)

func GetCocaWords() WordSet {
	return cocaSetsMemo.Get().(cocaSets).AllWords
}

func GetCocaAdjectives() WordSet {
	return cocaSetsMemo.Get().(cocaSets).Adjectives
}

type cocaSets struct {
	AllWords   WordSet
	Adjectives WordSet
}

var cocaSetsMemo = NewMemo(func() interface{} {
	reader := csv.NewReader(Get_coca_data("coca-5000.csv"))
	// Disable field count checking.
	reader.FieldsPerRecord = -1

	// Our raw data may contain 2 lines with the same word if that word can be
	// used as more than one part of speech. We just add the occurrence counts
	// of these lines together.
	var sets = &cocaSets{NewWordSet(), NewWordSet()}

	for i := 0; true; i++ {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}
		if i < 2 {
			// Skip the first 2 lines; they are headers.
			continue
		}
		if len(record) != 5 {
			log.Fatalln(
				"Wrong number of columns (", len(record), ") on line ", i)
		}

		occurrences, err := strconv.ParseInt(record[3], 10, 64)
		if err != nil {
			log.Fatalln("Invalid occurrences on line ", i)
		}
		word := record[1]
		partOfSpeech := record[2]

		if word == "n't" {
			word = "not"
		}
		adjective := strings.Index(partOfSpeech, "j") >= 0

		sets.AllWords.Add(WeightedWord{word, occurrences})
		if adjective {
			sets.Adjectives.Add(WeightedWord{word, occurrences})
		}
	}

	return sets
})
