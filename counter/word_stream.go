package counter

import (
  "bufio"
  "io"
  "strings"
  "unicode"
)

func ProcessWords(text io.Reader, process func(string) error) error {
  scanner := bufio.NewScanner(text)
  scanner.Split(bufio.ScanWords)
  for scanner.Scan() {
    var word string = scanner.Text()
    word = strings.TrimFunc(word, func(r rune) bool {
      return !unicode.IsLetter(r)
    })
    if len(word) == 0 {
      continue
    }
    word = strings.ToLower(word)
    err := process(word)
    if err != nil {
      return err
    }
  }
  return nil
}