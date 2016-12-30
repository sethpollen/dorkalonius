package gutenberg_test

import (
	"io"
	"log"
	"strings"
	"testing"
)
import . "github.com/sethpollen/dorkalonius/gutenberg"

const (
	goodEbook = `
Title: Adventures of Huckleberry Finn, Complete

*** START OF THIS PROJECT GUTENBERG EBOOK HUCKLEBERRY FINN ***

By Mark
Twain

*** END OF THIS PROJECT GUTENBERG EBOOK HUCKLEBERRY FINN ***

***** This file should be named 76-0.htm or 76-0.zip *****
`
	badEbook = `
Title: Adventures of Huckleberry Finn, Complete

By Mark
Twain
`
	maxBufSize = 2048
)

func TestBasic(t *testing.T) {
	for bufSize := 1; bufSize <= maxBufSize; bufSize++ {
		log.Printf("bufSize = %d", bufSize)
		reader := NewEbookReader(strings.NewReader(goodEbook))
		buf := make([]byte, maxBufSize)

		pos := 0
		for {
			endPos := maxBufSize
			if pos+bufSize < maxBufSize {
				endPos = pos + bufSize
			}
			n, err := reader.Read(buf[pos:endPos])
			pos += n
			if err == nil {
				if n == 0 {
					t.Error("Zero 'n' but no error")
					return
				}
				continue
			}
			if err == io.EOF {
				break
			}
			t.Error(err)
			return
		}

		result := strings.TrimSpace(string(buf[0:pos]))
		if result != "By Mark\nTwain" {
			t.Errorf("%q", result)
			return
		}
	}
}

func TestMalformed(t *testing.T) {
	reader := NewEbookReader(strings.NewReader(badEbook))
	buf := make([]byte, maxBufSize)
	n, err := reader.Read(buf)
	if n != 0 {
		t.Error(n)
	}
	if err != io.EOF {
		t.Error(err)
	}
}
