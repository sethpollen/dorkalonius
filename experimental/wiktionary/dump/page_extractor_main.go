// Utility to search over a Wiktionary dump and extract entire page contents.

package main

import (
	"flag"
	"fmt"
	"github.com/sethpollen/sbp_linux_config/sbpgo/games/words/wiktionary/dump"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

var inputFile = flag.String("input", "", "XML file to read")
var outputDir = flag.String("output_dir", "",
	"Directory to dump Module: page contents")

func main() {
	flag.Parse()
	if len(*inputFile) == 0 {
		log.Fatalln("--input is required")
	}
	if len(*outputDir) == 0 {
		log.Fatalln("--output_dir is required")
	}

	file, err := os.Open(*inputFile)
	if err != nil {
		log.Fatalln(err)
	}

	filesWritten := 0
	dump.ReadDump(file, func(page *dump.Page) {
		var interestingPrefix string
		for _, prefix := range []string{"Module", "Wiktionary"} {
			if strings.HasPrefix(page.Title, prefix+":") {
				interestingPrefix = prefix
				break
			}
		}

		if len(interestingPrefix) > 0 {
			// We can't use page.Title as the output filename, since page Titles
			// can contain characters like '/'.
			filename := fmt.Sprintf("%s_%d.txt", interestingPrefix, filesWritten)
			outputFile := path.Join(*outputDir, filename)
			filesWritten++

			// Include the page title at the beginning of the file.
			err = ioutil.WriteFile(outputFile,
				[]byte(fmt.Sprintf("%s\n\n%s", page.Title, page.Text)), 0666)
			if err != nil {
				log.Fatalln(err)
			}

			log.Printf("Extracted %s to %s", page.Title, filename)
		}
	})
}
