// Tool for converting CSV data into a serialized WordSet object.

package main

import (
	"encoding/csv"
	"flag"
	"github.com/sethpollen/dorkalonius"
	"io"
	"log"
	"os"
	"strconv"
)

// Input files are passed as plain command-line arguments.
var outputFile = flag.String("output_file", "",
	"Serialized WordSet file to write")

// CSV interpretation settings.
var csvHeaderLines = flag.Int("csv_header_lines", 0,
	"Number of CSV lines to skip at the beginning of the input file")
var csvWordColumn = flag.Int("csv_word_column", 0,
	"Column in the CSV file which contains the word")
var csvWeightColumn = flag.Int("csv_weight_column", 1,
	"Column in the CSV file which contains the weight (occurrences)")

func main() {
	flag.Parse()

	tasks := make([]func() dorkalonius.WordSet, flag.NArg())
	for i := range tasks {
		filename := flag.Arg(i)
		tasks[i] = func() dorkalonius.WordSet {
			return readFile(filename)
		}
	}
	wordSet := dorkalonius.BuildWordSet(tasks)

	out, err := os.Create(*outputFile)
	if err != nil {
		log.Fatal(err)
	}

	err = wordSet.Serialize(out)
	if err != nil {
		log.Fatal(err)
	}
}

func readFile(filename string) dorkalonius.WordSet {
	in, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	csvIn := csv.NewReader(in)
	// Disable field count checking.
	csvIn.FieldsPerRecord = -1

	wordSet := dorkalonius.NewWordSet()

	for i := 0; true; i++ {
		record, err := csvIn.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln(err)
		}
		if i < *csvHeaderLines {
			continue
		}

		word := record[*csvWordColumn]
		// TODO: word = strings.ToLower(word)
		weight, err := strconv.ParseInt(record[*csvWeightColumn], 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		wordSet.Add(dorkalonius.WeightedWord{word, weight})
	}

	return wordSet
}
