// Tool for converting the CSV wordlist file into a Go source file providing
// programmatic access to it without any runtime file dependencies.

package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/sethpollen/sbp_linux_config/sbpgo/games/words"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

var sourceFile = flag.String("source_file", "",
	"CSV file containing wordlist data")
var destFile = flag.String("dest_file", "",
	"Go file to write")

// Reads in the list of words from the file.
func ReadWordList(path string) (*words.WordList, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(file)
	// Disable field count checking.
	reader.FieldsPerRecord = -1

	// Our raw data may contain 2 lines with the same word if that word can be
	// used as more than one part of speech. We just add the occurrence counts
	// of these lines together.
	var wordSet = make(map[string]*words.Word)

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
			wordSet[word] = &words.Word{word, occurrences, partOfSpeech}
		}
	}

	// Convert the map to a WordList object.
	wordList := words.NewWordList()
	for _, word := range wordSet {
		wordList.AddWord(*word)
	}

	sort.Sort(wordList)
	return wordList, nil
}

func main() {
	flag.Parse()

	if *sourceFile == "" {
		log.Fatalln("--source_file is required")
	}
	if *destFile == "" {
		log.Fatalln("--dest_file is required")
	}

	list, err := ReadWordList(*sourceFile)
	if err != nil {
		log.Fatalln(err)
	}

	out, err := os.Create(*destFile)
	if err != nil {
		log.Fatalln(err)
	}

	var header = `
    package coca

    import "github.com/sethpollen/sbp_linux_config/sbpgo/games/words"

    func GetWordList() *words.WordList {
      return &words.WordList{[]words.Word{
    `
	var footer = fmt.Sprintf(`
      }, %d}
    }
    `, list.TotalOccurrences)

	out.Write([]byte(header))
	for _, word := range list.Words {
		out.Write([]byte(fmt.Sprintf("words.Word{%q, %d, %q},\n",
			word.Word, word.Occurrences, word.PartsOfSpeech)))
	}
	out.Write([]byte(footer))
}
