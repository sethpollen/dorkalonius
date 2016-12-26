package main

import (
	"flag"
	"fmt"
	"github.com/sethpollen/dorkalonius"
	"io"
	"log"
	"os"
)

var outputFile = flag.String("output_file", "", "Go source file to write")

func main() {
	flag.Parse()
	var err error

	out, err := os.Create(*outputFile)
	if err != nil {
		log.Fatalln(err)
	}

	_, err = out.WriteString(`
    package main
    
    import (
      "bytes"
      "compress/flate"
      "encoding/base64"
      "io"
      "strings"
      "testing"
    )

    func Test(t *testing.T) {
      var expected []byte
      var reader io.Reader
      var actual []byte
      var n int
      var err error
  `)
	if err != nil {
		log.Fatalln(err)
	}

	emitTestCase(out, make([]byte, 0))
	emitTestCase(out, []byte("ABC"))

	data := make([]byte, 256)
	for i := 0; i < 256; i++ {
		data[i] = byte(i)
	}
	emitTestCase(out, data)

	_, err = out.WriteString("}")
	if err != nil {
		log.Fatalln(err)
	}
}

func emitTestCase(out io.Writer, data []byte) {
	var err error

	_, err = out.Write([]byte(fmt.Sprintf(`
    expected = []byte(%q)
    actual = make([]byte, %d)
    reader =
    `, data, len(data))))
	if err != nil {
		log.Fatalln(err)
	}

	encoder, err := dorkalonius.NewGoEmbedEncoder(out)
	if err != nil {
		log.Fatalln(err)
	}
	encoder.Write(data)
	encoder.Close()

	_, err = out.Write([]byte(`
    n, err = io.ReadFull(reader, actual)
    if err != nil {
      t.Errorf("Generated reader produced error: %v (read %d bytes)", err, n)
    } else {
      if bytes.Compare(expected, actual[0:n]) != 0 {
        t.Errorf("Expected: %q\nActual: %q", expected, actual[0:n])
      }
      n, _ = io.ReadFull(reader, make([]byte, 1))
      if n > 0 {
        t.Error("Found unexpected bytes")
      }
    }
  `))
	if err != nil {
		log.Fatalln(err)
	}
}
