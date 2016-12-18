// Tool for converting the inflections.xml file into a Go source file providing
// programmatic access to it without any runtime file dependencies.

package main

import (
  "bytes"
  "encoding/base64"
  "encoding/gob"
  "encoding/xml"
	"flag"
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

	in, err := os.Open(*sourceFile)
  if err != nil {
    log.Fatalln(err)
  }
  
  xmlDecoder := xml.NewDecoder(in)
  var inflections wiktionary.Inflections
  err = xmlDecoder.Decode(&inflections)
  if err != nil {
    log.Fatalln(err)
  }

  // Sanity check that we didn't skip a bunch of XML.
  if len(inflections.Inflections) == 0 {
    log.Fatalln("No inflections read from input file")
  }
  
  inflectionMap := wiktionary.NewInflectionMap()
  for _, inflection := range inflections.Inflections {
    for _, inflectedForm := range inflection.InflectedForms {
      inflectionMap.Add(inflection.BaseWord, inflectedForm)
    }
  }
  
  // TODO: This causes bazel to hang, so we skip it for now.
  log.Fatalln("TODO:")
  
  var encodedMap bytes.Buffer
  base64Encoder := base64.NewEncoder(base64.StdEncoding, &encodedMap)
  gobEncoder := gob.NewEncoder(base64Encoder)
  if err = gobEncoder.Encode(inflectionMap); err != nil {
    log.Fatalln(err)
  }
  base64Encoder.Close()
  
	out, err := os.Create(*destFile)
	if err != nil {
		log.Fatalln(err)
	}

	out.Write([]byte(`
    package embed
    
    import (
      "encoding/base64"
      "encoding/gob"
      "github.com/sethpollen/dorkalonius/incubator/wiktionary"
      "strings"
    )

    func GetInflectionMap() *wiktionary.InflectionMap {
      reader := strings.NewReader(encodedList)
      base64Decoder := base64.NewDecoder(base64.StdEncoding, reader)
      gobDecoder := gob.NewDecoder(base64Decoder)
      var inflectionMap wiktionary.InflectionMap
      gobDecoder.Decode(&inflectionMap)
      return &inflectionMap
    }

    const encodedList =
    `))

  for encodedMap.Len() > 0 {
    out.Write([]byte("\""))
    out.Write(encodedMap.Next(75))
    out.Write([]byte("\"+\n"))
  }
  // Close the final + sign with an empty string.
  out.Write([]byte("\"\"\n"))
}
