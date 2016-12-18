package dorkalonius

import (
  "fmt"
  "io"
)

// Returns a WriteCloser which will encode its bytes as a Go string constant
// and write the result to 'sink'. The returned object must be Closed() to
// finish writing a complete Go string constant.
func NewGoEmbedEncoder(sink io.Writer) io.WriteCloser {
  sink.Write([]byte("\n\""))
  return encoder{sink, 0}
}

// Max number of data characters to write to the current line. Doesn't account
// for the potential overhead of escaping.
const maxLineSize = 80

type encoder struct {
  sink io.Writer
  // Number of data characters written to the current line.
  lineSize int
}

func (self encoder) Write(p []byte) (int, error) {
  bytesWritten := 0
  for len(p) > 0 {
    n := maxLineSize - self.lineSize
    if len(p) < n {
      n = len(p)
    }
    data := p[0:n]
    p = p[n:]

    // Apply escaping.
    data = []byte(fmt.Sprintf("%q", data))
    // Remove double-quotes added by Sprintf.
    data = data[1:len(data)-1]
    
    if _, err := self.sink.Write(data); err != nil {
      return bytesWritten, err
    }
    self.lineSize += n
    bytesWritten += n
    
    if self.lineSize >= maxLineSize {
      if _, err := self.sink.Write([]byte("\"+\n\"")); err != nil {
        return bytesWritten, err
      }
      self.lineSize = 0
    }
  }
  return bytesWritten, nil
}

func (self encoder) Close() error {
  self.sink.Write([]byte("\"\n"))
  return nil
}