// Program to embed a file as data in Go source code. The output is a .go
// file which declares a single method. This method returns an io.Reader which
// yields the embedded data.

package main

import (
	"flag"
	"fmt"
	"github.com/sethpollen/dorkalonius/tools"
	"io"
	"log"
	"os"
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
	if flag.NArg() != 1 {
    log.Fatalln("need exactly 1 input file")
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
      "strings"
    )
  
    func %s() io.Reader {
      return `, *packageName, *methodName))
	if err != nil {
		log.Fatalln(err)
	}

  in, err := os.Open(flag.Arg(0))
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

	if _, err = out.WriteString("}\n"); err != nil {
		log.Fatalln(err)
	}
}
