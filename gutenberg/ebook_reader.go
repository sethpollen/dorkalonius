// Provides an io.Reader for extracting the actuall ebook text from a plaintext
// file from gutenberg.org.

package gutenberg

import (
	"bufio"
	"bytes"
	"io"
)

func NewEbookReader(source io.Reader) io.Reader {
	return &reader{bufio.NewScanner(source), header, nil, false}
}

// Possible reader states.
const (
	header = iota
	body
	footer
)

type reader struct {
	Source      *bufio.Scanner
	State       int
	Buffer      []byte
	EmitNewline bool
}

const (
	prefixEndTag   = "*** START OF THIS PROJECT GUTENBERG EBOOK"
	suffixBeginTag = "*** END OF THIS PROJECT GUTENBERG EBOOK"
)

func (self *reader) getSourceError() error {
	err := self.Source.Err()
	if err == nil {
		return io.EOF
	}
	return err
}

func (self *reader) Read(p []byte) (int, error) {
	var n int = 0
	for n < len(p) {
		bytesWanted := len(p) - n

		if self.EmitNewline {
			p[n] = '\n'
			n++
			self.EmitNewline = false
			continue
		}

		switch self.State {
		case header:
			if !self.Source.Scan() {
				return n, self.getSourceError()
			}
			if bytes.HasPrefix(self.Source.Bytes(), []byte(prefixEndTag)) {
				self.State = body
				break
			}

		case body:
			if self.Buffer != nil && len(self.Buffer) > 0 {
				if len(self.Buffer) < bytesWanted {
					bytesWanted = len(self.Buffer)
				}
				copy(p[n:n+bytesWanted], self.Buffer[0:bytesWanted])
				n += bytesWanted
				self.Buffer = self.Buffer[bytesWanted:]
				if len(self.Buffer) == 0 {
					// Add back in the newline eaten by the Scanner.
					self.EmitNewline = true
				}
				break
			}
			if !self.Source.Scan() {
				return n, self.getSourceError()
			}
			buffer := self.Source.Bytes()
			if bytes.HasPrefix(buffer, []byte(suffixBeginTag)) {
				self.State = footer
				break
			}
			self.Buffer = buffer

		case footer:
			return n, io.EOF
		}
	}
	return n, nil
}
