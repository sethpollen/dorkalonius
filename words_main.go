package main

import (
	"flag"
	"fmt"
	"github.com/sethpollen/dorkalonius"
	"log"
	"math/rand"
	"os"
	"path"
	"strconv"
	"time"
)

var sample_size = flag.Int("sample_size", 35,
	"Number of words to sample.")
var outputWidth = flag.Int("output_width", -1,
	"Width of the terminal where output will be shown.")
var duration = flag.Duration("duration", 0,
	"Duration for the game timer which runs after words are printed.")

var outputDir = flag.String("output_dir", "",
	"If provided, output will be written to files in this directory instead "+
		"of to stdout.")
var outputFiles = flag.Int("output_files", 1,
	"Only used if --output_dir is provided. Gives the number of unique game "+
		"files to generate. Each file will have the format N.txt, where N is "+
		"an integer (possibly zero-padded) between 0 and --output_files.")

func main() {
	flag.Parse()
	rand.Seed(time.Now().UTC().UnixNano())

	if *sample_size < 0 {
		log.Fatalln("--sample_size must be nonnegative")
	}
	var err error

	sampler, err := dorkalonius.GetCocaIndex()
	if err != nil {
		log.Fatalln(err)
	}

	if *outputDir == "" {
		fmt.Println()
		err = generateGame(sampler, os.Stdout)
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Println()

		if *duration > 0 {
			dorkalonius.VerboseSleep(*duration, true)
			fmt.Println("TIME'S UP")
			fmt.Println()
		}
		return
	}

	if *outputFiles <= 0 {
		log.Fatalln("--output_files must be positive")
	}
	fileNameFormat := fmt.Sprintf("words_%d_%%0%dd.txt",
		*sample_size,
		len(strconv.Itoa(*outputFiles-1)))

	for i := 0; i < *outputFiles; i++ {
		out, err := os.Create(path.Join(*outputDir,
			fmt.Sprintf(fileNameFormat, i)))
		if err != nil {
			log.Fatalln(err)
		}
		err = generateGame(sampler, out)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func generateGame(sampler *dorkalonius.Index, out *os.File) error {
	var err error

	game, err := dorkalonius.NewGame(sampler)
	if err != nil {
		return err
	}

	_, err = out.WriteString(fmt.Sprintf("TARGET WORD: %s\n\n",
		game.TargetWord))
	if err != nil {
		return err
	}
	_, err = out.WriteString("AVAILABLE WORDS:\n\n")
	if err != nil {
		return err
	}
	err = printWords(game.AvailableWords, out)
	if err != nil {
		return err
	}
	return nil
}

// Pretty-print words in columns to 'out'.
func printWords(wordList *dorkalonius.WordList, out *os.File) error {
	screenWidth := *outputWidth
	if screenWidth < 1 {
		screenWidth = 1
	}

	// We take a simple approach by using the same width for all columns. Find
	// the longest word to determine that width.
	var maxWordLength int = 0
	for _, word := range wordList.Words {
		if len(word.Word) > maxWordLength {
			maxWordLength = len(word.Word)
		}
	}
	columnWidth := maxWordLength + 3

	columns := int(screenWidth) / columnWidth
	if columns < 1 {
		columns = 1
	}
	rows := (wordList.Len() + columns - 1) / columns

	// We print down each column, then across.
	var err error
	for row := 0; row < rows; row++ {
		for col := 0; col < columns; col++ {
			var index = row + (col * rows)
			if index >= wordList.Len() {
				continue
			}
			_, err = out.WriteString(wordList.Words[index].Word)
			if err != nil {
				return err
			}
			for i := 0; i < columnWidth-len(wordList.Words[index].Word); i++ {
				_, err = out.WriteString(" ")
				if err != nil {
					return err
				}
			}
		}
		_, err = out.WriteString("\n")
		if err != nil {
			return err
		}
	}

	return nil
}
