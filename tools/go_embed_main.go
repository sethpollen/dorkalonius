// Program to embed a file as data in Go source code. The output is a .go
// file which declares a single method. This method accepts a filename and
// an io.Reader which yields the embedded data for the given filename.

package main

import (
	"flag"
	"fmt"
	"github.com/sethpollen/dorkalonius/tools"
	"io"
	"log"
	"os"
	"path"
)

var outputFile = flag.String("output_file", "", "Go source file to write")
var packageName = flag.String("package", "", "Go package name to use")
var methodName = flag.String("method", "", "Go method name to use")

// Remaining command-line arguments are input files.

func main() {
	flag.Parse()
	var err error

	if len(*outputFile) == 0 {
		log.Fatalln("missing --output_file")
	}
	if len(*packageName) == 0 {
		log.Fatalln("missing --package")
	}
	if len(*methodName) == 0 {
		log.Fatalln("missing --method")
	}

	out, err := os.Create(*outputFile)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = out.WriteString(fmt.Sprintf(`
    package %s
    
    import (
      "compress/flate"
      "encoding/base64"
      "io"
      "log"
      "strings"
    )
  
    func %s(file string) io.Reader {
      switch file {
    `, *packageName, *methodName))
	if err != nil {
		log.Fatalln(err)
	}

	for _, inputFile := range flag.Args() {

		_, err = out.WriteString(fmt.Sprintf(`
          case %q:
            return `, path.Base(inputFile)))
		if err != nil {
			log.Fatalln(err)
		}

		in, err := os.Open(inputFile)
		if err != nil {
			log.Fatalln(err)
		}

		encoder, err := tools.NewGoEmbedEncoder(out)
		if err != nil {
			log.Fatalln(err)
		}
		if _, err = io.Copy(encoder, in); err != nil {
			log.Fatalln(err)
		}
		if err = encoder.Close(); err != nil {
			log.Fatalln(err)
		}
	}

	_, err = out.WriteString(`
     }
     log.Fatalf("Unrecognized filename: %q", file)
     return nil
   }
  `)
	if err != nil {
		log.Fatalln(err)
	}
}
