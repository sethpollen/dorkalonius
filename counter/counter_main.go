// Tool for digesting corpora and producing word counts. Reads text from
// stdin and then prints a report.

package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/sethpollen/dorkalonius"
	"github.com/sethpollen/dorkalonius/counter"
	"github.com/sethpollen/dorkalonius/gutenberg"
	"github.com/sethpollen/dorkalonius/wiktionary"
	"io"
	"log"
	"os"
	"sort"
)

var gutenbergEbook = flag.Bool("gutenberg_ebook", false,
	"If true, interpret input files as Project Gutenberg ebooks.")

// Accepts a list of input files as command-line arguments.
func main() {
	flag.Parse()

	inflectionMap, err := wiktionary.InflectionMapFromBzippedXml(
		"./wiktionary/inflections.xml.bz2")
	if err != nil {
		log.Fatalln(err)
	}

	wordChan := make(chan string, 1000)

	for _, filename := range flag.Args() {
		go worker(filename, wordChan, inflectionMap)
	}

	// Collect outputs from workers.
	wordCountMap := make(map[string]int64)
	workersDone := 0
	for workersDone < flag.NArg() {
		word := <-wordChan
		if len(word) == 0 {
			// This is the sentinel value sent by a worker as it exits.
			workersDone++
			continue
		}
		wordCountMap[word]++
	}

	wordList := dorkalonius.NewWordList()
	for wordStr, occurrences := range wordCountMap {
		if occurrences <= 1 {
			// We drop words which occur only once in the whole corpus. These are
			// likely not to be real words at all, but rather misrecognized patterns
			// like "foo--bar".
			//
			// TODO: Make some more clever logic to drop any words with "--" in the
			// mindle.
			continue
		}
		wordList.AddWord(dorkalonius.Word{wordStr, occurrences, false})
	}
	sort.Sort(wordList)

	csvWriter := csv.NewWriter(os.Stdout)
	for _, word := range wordList.Words {
		csvWriter.Write([]string{word.Word, fmt.Sprintf("%d", word.Occurrences)})
	}
	csvWriter.Flush()
}

func worker(filename string, wordChan chan string,
	inflectionMap *wiktionary.InflectionMap) {
	var input io.Reader
	var err error
	input, err = os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}

	if *gutenbergEbook {
		input = gutenberg.NewEbookReader(input)
	}

	counter.ProcessWords(input, func(word string) error {
		word = inflectionMap.GetBaseWord(word)
		if len(word) == 0 {
			log.Fatalln("Empty word")
		}
		wordChan <- word
		return nil
	})
	// Send a sentinel value to indicate that this worker is finished.
	wordChan <- ""
}
