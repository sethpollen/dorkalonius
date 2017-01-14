// Tool for converting CSV data into a serialized WordSet object.

package main

import (
	"encoding/csv"
	"flag"
	"github.com/sethpollen/dorkalonius/util"
	"io"
	"log"
	"os"
	"strconv"
  "strings"
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
var csvFilterColumn = flag.Int("csv_filter_column", -1,
  "Column in the CSV file which contains cell to filter on. Leave "+
  "absent to specify no filtering")
var csvFilterValue = flag.String("csv_filter_value", "",
  "We only keep rows where the csv_filter_column has this value")

func main() {
	flag.Parse()

	tasks := make([]func() util.WordSet, flag.NArg())
	for i := range tasks {
		filename := flag.Arg(i)
		tasks[i] = func() util.WordSet {
			return readFile(filename)
		}
	}
	wordSet := util.BuildWordSet(tasks)

	out, err := os.Create(*outputFile)
	if err != nil {
		log.Fatal(err)
	}

	err = wordSet.Serialize(out)
	if err != nil {
		log.Fatal(err)
	}
}

func readFile(filename string) util.WordSet {
	in, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	csvIn := csv.NewReader(in)
	// Disable field count checking.
	csvIn.FieldsPerRecord = -1

	wordSet := util.NewWordSet()

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
		
		if *csvFilterColumn >= 0 {
      if record[*csvFilterColumn] != *csvFilterValue {
        continue
      }
    }

		word := record[*csvWordColumn]
		word = strings.ToLower(word)
		weight, err := strconv.ParseInt(record[*csvWeightColumn], 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		wordSet.Add(util.WeightedWord{word, weight})
	}

	return wordSet
}
