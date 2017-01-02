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

	responseChans := make([]chan dorkalonius.WordSet, flag.NArg())
	for i, filename := range flag.Args() {
    responseChans[i] = make(chan dorkalonius.WordSet)
		go worker(inflectionMap, filename, responseChans[i])
	}

	// Collect outputs from workers.
	wordSet := dorkalonius.NewWordSet()
	for _, responseChan := range responseChans {
    wordSet.AddAll(<-responseChan)
	}

	csvWriter := csv.NewWriter(os.Stdout)
	for _, word := range wordSet.GetWords() {
		csvWriter.Write([]string{word.Word, fmt.Sprintf("%d", word.Weight)})
	}
	csvWriter.Flush()
}

func worker(
  inflectionMap *wiktionary.InflectionMap,
  filename string,
  responseChan chan<- dorkalonius.WordSet) {

	var input io.Reader
	var err error
	input, err = os.Open(filename)
	if err != nil {
		log.Fatalln(err)
	}

	if *gutenbergEbook {
		input = gutenberg.NewEbookReader(input)
	}

	wordSet := dorkalonius.NewWordSet()
	err = counter.ProcessWords(input, func(word string) error {
		word = inflectionMap.GetBaseWord(word)
		if len(word) == 0 {
			log.Fatalln("Empty word")
		}
		wordSet.Add(dorkalonius.WeightedWord{word, 1})
		return nil
	})
  if err != nil {
    log.Fatalln(err)
  }

  responseChan <- wordSet
}
