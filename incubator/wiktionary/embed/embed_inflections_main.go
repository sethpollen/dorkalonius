// Tool for converting the inflections.xml file into a Go source file providing
// programmatic access to it without any runtime file dependencies.

package main

import (
  "encoding/xml"
	"flag"
  "fmt"
	"github.com/sethpollen/dorkalonius/incubator/wiktionary"
	"log"
	"os"
)

var sourceFile = flag.String("source_file", "",
	"XML file containing inflection data")
var destFile = flag.String("dest_file", "",
	"Go file to write")

func main() {
	flag.Parse()

	if *sourceFile == "" {
		log.Fatalln("--source_file is required")
	}
	if *destFile == "" {
		log.Fatalln("--dest_file is required")
	}

	inFile, err := os.Open(*sourceFile)
  if err != nil {
    log.Fatalln(err)
  }
  
  decoder := xml.NewDecoder(inFile)
  var inflections wiktionary.Inflections
  err = decoder.Decode(&inflections)
  if err != nil {
    log.Fatalln(err)
  }

  if len(inflections.Inflections) == 0 {
    log.Fatalln("No inflections read from input file")
  }
  
	outFile, err := os.Create(*destFile)
	if err != nil {
		log.Fatalln(err)
	}

	var header = `
    package embed
    
    import "github.com/sethpollen/dorkalonius/incubator/wiktionary"

    func GetInflectionMap() *wiktionary.InflectionMap {
      i := wiktionary.NewInflectionMap()
    `
	var footer = `
	    return i
    }
    `

	outFile.Write([]byte(header))
  
	for _, inflection := range inflections.Inflections {
		for _, inflectedForm := range inflection.InflectedForms {
      outFile.Write([]byte(fmt.Sprintf("i.Add(%q, %q)\n",
        inflection.BaseWord, inflectedForm)))
    }
	}
	
	outFile.Write([]byte(footer))
}
