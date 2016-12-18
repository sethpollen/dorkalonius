package dorkalonius_test

import (
  "bytes"
  "math/rand"
  "testing"
)
import . "github.com/sethpollen/dorkalonius"

func encode(data string) string {
  var buffer bytes.Buffer
  encoder := NewGoEmbedEncoder(&buffer)
  
  dataBytes := []byte(data)
  rand.Seed(int64(len(dataBytes)))
  for len(dataBytes) > 0 {
    n := 0
    if rand.Intn(3) > 0 {
      n = rand.Intn(len(dataBytes) + 1)
    }
    encoder.Write(dataBytes[0:n])
    dataBytes = dataBytes[n:]
  }

  encoder.Close()
  return buffer.String()
}

func TestEmpty(t *testing.T) {
  actual := encode("")
  if actual != "\n\"\"\n" {
    t.Errorf("%q", actual)
  }
}

func TestOneChar(t *testing.T) {
  actual := encode("A")
  if actual != "\n\"A\"\n" {
    t.Errorf("%q", actual)
  }
}

func TestOneCharEscape(t *testing.T) {
  actual := encode("\\")
  if actual != "\n\"\\\\\"\n" {
    t.Errorf("%q", actual)
  }
}

func TestMultipleLines(t *testing.T) {
  var buffer bytes.Buffer
  for i := 0; i < 100; i++ {
    buffer.Write([]byte("A"))
  }
  
  actual := encode(buffer.String())
  if actual != "\n\"AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"+
    "AAAAAAAAAAAAAAAAAAAAAAA\"+\n\"AAAAAAAAAAAAAAAAAAAA\"\n" {
    t.Errorf("%q", actual)
  }
}
