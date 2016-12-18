package main

import (
	"encoding/csv"
	"flag"
	"github.com/sethpollen/sbp_linux_config/sbpgo/games/words/wiktionary"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

var skipLines = flag.Int("skip_lines", 0, "Initial CSV lines to be skipped.")

const inputCsvFile = "./sbpgo/games/words/wiktionary/dump/data/en-templates.csv"
const concurrency = 16

// Argument/return types for the ExpandInflections call.
type InflectionRequest struct {
	Line      int
	CsvRecord []string
}
type InflectionResponse struct {
	Line        int
	Pos         int
	Title       string
	Inflections []string
	// Only set to true for highly irregular words. This is a hint that people
	// might need to be told explicitly what the inflections are.
	Irregular   bool
}

// Method run by worker threads. Will send nil to 'responseChan' once it
// finishes processing everything from 'requestChan'.
func worker(requestChan <-chan InflectionRequest,
	responseChan chan<- *InflectionResponse) {
	inflector, err := wiktionary.NewInflector()
	if err != nil {
		log.Fatalln(err)
	}

	// Regexes to match MediaWiki-style link expressions, such as the following:
	//  [[link]]        (simple)
	//  [[link|title]]  (aliased)
	simpleLinkRe := regexp.MustCompile("\\[\\[([^\\|\\]]+)\\]\\]")
	aliasedLinkRe := regexp.MustCompile("\\[\\[[^\\|\\]]+\\|([^\\|\\]]+)\\]\\]")

	for request := range requestChan {
		var title string = request.CsvRecord[0]
		if strings.Index(title, " ") >= 0 {
			// Drop any multi-word forms, as we only process corpora one word at a
			// time.
			continue
		}

		var invocation string = request.CsvRecord[1]
		
    if strings.Index(invocation, "highly irregular") >= 0 {
      log.Fatalf("Cannot handle highly irregular entry: %q", invocation)
    }
		
		invocation = strings.TrimPrefix(invocation, "{{en-")
		invocation = strings.TrimSuffix(invocation, "}}")
		invocation = simpleLinkRe.ReplaceAllString(invocation, "$1")
		invocation = aliasedLinkRe.ReplaceAllString(invocation, "$1")

		var invocationParts []string = strings.Split(invocation, "|")

		var partOfSpeech = invocationParts[0]
		var posEnum int
		switch partOfSpeech {
		case "noun":
			posEnum = wiktionary.Noun
		case "verb":
			posEnum = wiktionary.Verb
		case "adj":
			posEnum = wiktionary.Adjective
		case "adv", "adverb":
			posEnum = wiktionary.Adverb
		default:
			log.Fatalf("Unrecognized part of speech on line %d: %s",
				request.Line, partOfSpeech)
		}

		var args []string = invocationParts[1:]
		expanded, err := inflector.ExpandInflections(
      posEnum, title, args)
		if err != nil {
			log.Fatalf("Inflector failed on CSV line %d:\n%s\n%s",
                 request.Line, strings.Join(request.CsvRecord, ", "), err)
		}
		
		if len(expanded) == 0 {
      // There are no inflections to converge to the base form.
      continue
    }

		responseChan <- &InflectionResponse{
      request.Line, posEnum, title, expanded, false}
	}
	responseChan <- nil
}

func main() {
	flag.Parse()

	file, err := os.Open(inputCsvFile)
	if err != nil {
		log.Fatalln(err)
	}
	csv := csv.NewReader(file)

	// We farm out the Lua invocations to several goroutines for parallelism.
	requestChan := make(chan InflectionRequest, 100)
	responseChan := make(chan *InflectionResponse, 100)
	for i := 0; i < concurrency; i++ {
		go worker(requestChan, responseChan)
	}
	
	// Manually insert an entry for the verb "be". This is the only page on the
	// English Wiktionary that invokes the "highly irregular" cop-out.
	responseChan <- &InflectionResponse{
    0, wiktionary.Verb, "be",
    []string{"am", "is", "are", "was", "were", "being", "beings", "been"},
    true}

	// Spawn another goroutine to read in the CSV file and distribute its lines
	// to the workers.
	go func() {
		for line := 1; ; line++ {
			record, err := csv.Read()
			if err == io.EOF {
				break
			} else if err != nil {
				log.Fatalln(err)
			}
			if line <= *skipLines {
				continue
			}
			requestChan <- InflectionRequest{line, record}
		}
		close(requestChan)
	}()
  
  type BaseWord struct {
    Word string
    Irregular bool
  }

	// Collect and print the results in the main thread.
	inflectionToBaseWord := make(map[string]*BaseWord)

	// Count the number of nils; this indicates how many workers have
	// completed.
	nils := 0
	for nils < concurrency {
		response := <-responseChan
		if response == nil {
			nils++
			continue
		}
		
		baseWord := BaseWord{response.Title, response.Irregular}
		for _, inflection := range response.Inflections {
      existingBaseWord, ok := inflectionToBaseWord[inflection]
      if ok {
        if len(existingBaseWord.Word) <= len(baseWord.Word) {
          // We prefer the existing base word, as it's shorter.
          continue
        }
      }
      inflectionToBaseWord[inflection] = &baseWord
    }

		// TODO:
	}
}
