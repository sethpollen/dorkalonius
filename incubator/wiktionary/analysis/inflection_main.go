package main

import (
	"encoding/csv"
	"encoding/xml"
	"flag"
  "github.com/sethpollen/dorkalonius/incubator/wiktionary"
  "github.com/sethpollen/dorkalonius/incubator/wiktionary/analysis"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

var skipLines = flag.Int("skip_lines", 0,
	"Initial CSV lines to be skipped from input file.")
var outputXmlFile = flag.String("output_file", "",
	"Output XML file to write. See inflection_xml.go for the format.")

const inputCsvFile = "./incubator/wiktionary/dump/data/en-templates.csv"
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
}

// Method run by worker threads. Will send nil to 'responseChan' once it
// finishes processing everything from 'requestChan'.
func worker(requestChan <-chan InflectionRequest,
	responseChan chan<- *InflectionResponse) {
	inflector, err := analysis.NewInflector()
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
			posEnum = analysis.Noun
		case "verb":
			posEnum = analysis.Verb
		case "adj":
			posEnum = analysis.Adjective
		case "adv", "adverb":
			posEnum = analysis.Adverb
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

		// Note that 'expanded' may be empty if the base word has no other
		// forms.
		responseChan <- &InflectionResponse{
			request.Line, posEnum, title, expanded}
	}
	responseChan <- nil
}

func main() {
	flag.Parse()
	if len(*outputXmlFile) == 0 {
		log.Fatalln("--output_file is required")
	}

	inFile, err := os.Open(inputCsvFile)
	if err != nil {
		log.Fatalln(err)
	}
	inCsv := csv.NewReader(inFile)

	outFile, err := os.Create(*outputXmlFile)
	if err != nil {
		log.Fatalln(err)
	}

	// We farm out the Lua invocations to several goroutines for parallelism.
	requestChan := make(chan InflectionRequest, 100)
	responseChan := make(chan *InflectionResponse, 100)
	for i := 0; i < concurrency; i++ {
		go worker(requestChan, responseChan)
	}

	// Manually insert an entry for the verb "be". This is the only page on the
	// English Wiktionary that invokes the "highly irregular" cop-out.
	responseChan <- &InflectionResponse{
		0, analysis.Verb, "be",
		[]string{"am", "is", "are", "was", "were", "being", "beings", "been"}}

	// Spawn another goroutine to read in the CSV file and distribute its lines
	// to the workers.
	go func() {
		for line := 1; ; line++ {
			record, err := inCsv.Read()
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

	if _, err = outFile.Write([]byte("<inflections>\n")); err != nil {
		log.Fatalln(err)
	}

	// Collect and output the results in the main thread. Count the number of
	// nils; this indicates how many workers have completed.
	records := 0
	nils := 0
	for nils < concurrency {
		response := <-responseChan
		if response == nil {
			nils++
			continue
		}

		record := wiktionary.Inflection{
			BaseWord:       response.Title,
			Pos:            analysis.PosName(response.Pos),
			InflectedForms: response.Inflections,
		}
		xmlRecord, err := xml.MarshalIndent(record, "  ", "  ")
		if err != nil {
			log.Fatalln(err)
		}
		if _, err = outFile.Write(xmlRecord); err != nil {
			log.Fatalln(err)
		}
		if _, err = outFile.Write([]byte("\n")); err != nil {
			log.Fatalln(err)
		}

		records++
		if records%1000 == 0 {
			log.Printf("Processed %v records\n", records)
		}
	}

	if _, err = outFile.Write([]byte("</inflections>\n")); err != nil {
		log.Fatalln(err)
	}
}
