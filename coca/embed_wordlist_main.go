// Tool for converting the CSV wordlist file into a Go source file providing
// programmatic access to it without any runtime file dependencies.

package main

import (
  "bytes"
  "encoding/base64"
  "encoding/csv"
  "encoding/gob"
	"flag"
	"fmt"
	"github.com/sethpollen/dorkalonius"
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
func readWordList(path string) (*dorkalonius.WordList, error) {
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

func main() {
	flag.Parse()

	if *sourceFile == "" {
		log.Fatalln("--source_file is required")
	}
	if *destFile == "" {
		log.Fatalln("--dest_file is required")
	}

	list, err := readWordList(*sourceFile)
	if err != nil {
		log.Fatalln(err)
	}
	
	var encodedList bytes.Buffer
	base64Encoder := base64.NewEncoder(base64.StdEncoding, &encodedList)
	gobEncoder := gob.NewEncoder(base64Encoder)
  if err = gobEncoder.Encode(list); err != nil {
    log.Fatalln(err)
  }
  base64Encoder.Close()

	out, err := os.Create(*destFile)
	if err != nil {
		log.Fatalln(err)
	}
	
	out.Write([]byte(`
	  package coca
	  
	  import (
      "encoding/base64"
      "encoding/gob"
      "github.com/sethpollen/dorkalonius"
      "strings"
    )
  
    func GetWordList() *dorkalonius.WordList {
      reader := strings.NewReader(encodedList)
      base64Decoder := base64.NewDecoder(base64.StdEncoding, reader)
      gobDecoder := gob.NewDecoder(base64Decoder)
      var list dorkalonius.WordList
      gobDecoder.Decode(&list)
      return &list
    }
  
    const encodedList =
	`))
  
  for encodedList.Len() > 0 {
    out.Write([]byte("\""))
    out.Write(encodedList.Next(75))
    out.Write([]byte("\"+\n"))
  }
  // Close the final + sign with an empty string.
  out.Write([]byte("\"\"\n"))
}
