package tools

import (
	"bytes"
	"compress/flate"
	"encoding/base64"
	"io"
)

// Returns a WriteCloser which will encode its bytes as a Go expression and pass
// the source for this expression to 'sink'. The expression is guaranteed to
// produce an io.ReadCloser object which returns the encoded bytes. This expression
// will assume that the following imports are declared:
//   import (
//     "compress/flate"
//     "encoding/base64"
//     "strings"
//   )
func NewGoEmbedEncoder(sink io.Writer) (io.WriteCloser, error) {
	sink.Write([]byte("flate.NewReader(\n" +
		"base64.NewDecoder(base64.StdEncoding,\n" +
		"strings.NewReader(\n" +
		"\""))

	var dataBuffer bytes.Buffer

	compressor1 := base64.NewEncoder(base64.StdEncoding, &dataBuffer)
	compressor2, err := flate.NewWriter(compressor1, 2)
	if err != nil {
		return nil, err
	}

	return &encoder{[]io.WriteCloser{compressor2, compressor1},
		&dataBuffer, sink, 0}, nil
}

// Max number of data characters to write to the current line. Doesn't account
// for the potential overhead of escaping.
const maxLineSize = 100

type encoder struct {
	// Chain of steps which eventually dump into 'dataBuffer'.
	compressorChain []io.WriteCloser

	// Contains base64-encoded data.
	dataBuffer *bytes.Buffer

	// Sink for the final Go code.
	goCodeSink io.Writer

	// Number of data characters written to the current line.
	lineSize int
}

func (self *encoder) flush() error {
	for self.dataBuffer.Len() > 0 {
		data := make([]byte, maxLineSize-self.lineSize)
		n, err := self.dataBuffer.Read(data)
		if err != nil {
			return err
		}
		data = data[0:n]

		if _, err := self.goCodeSink.Write(data); err != nil {
			return err
		}
		self.lineSize += n

		if self.lineSize >= maxLineSize {
			if _, err := self.goCodeSink.Write([]byte("\"+\n\"")); err != nil {
				return err
			}
			self.lineSize = 0
		}
	}
	return nil
}

func (self *encoder) Write(p []byte) (int, error) {
	bytesAccepted, err := self.compressorChain[0].Write(p)
	if err != nil {
		return bytesAccepted, err
	}
	if err = self.flush(); err != nil {
		return bytesAccepted, err
	}
	return bytesAccepted, nil
}

func (self *encoder) Close() error {
	for _, compressor := range self.compressorChain {
		compressor.Close()
	}
	if err := self.flush(); err != nil {
		return err
	}
	self.goCodeSink.Write([]byte("\")))\n"))
	return nil
}
