// Tool for digesting corpora and producing word counts.

package main

import (
  "bufio"
  "fmt"
  "github.com/sethpollen/dorkalonius"
  "github.com/sethpollen/dorkalonius/counter"
  "github.com/sethpollen/dorkalonius/wiktionary"
  "log"
  "os"
  "sort"
)

func main() {
  inflectionMap, err := wiktionary.InflectionMapFromBzippedXml(
    "./wiktionary/inflections.xml.bz2")
  if err != nil {
    log.Fatalln(err)
  }
  
  wordCountMap := make(map[string]int64)
  input := bufio.NewReader(os.Stdin)
  counter.ProcessWords(input, func(word string) error {
    word = inflectionMap.GetBaseWord(word)
    wordCountMap[word]++
    return nil
  })
  
  wordList := dorkalonius.NewWordList()
  for wordStr, occurrences := range wordCountMap {
    wordList.AddWord(dorkalonius.Word{wordStr, occurrences, false})
  }
  sort.Sort(wordList)
  
  for _, word := range wordList.Words {
    fmt.Printf("%6d  %s\n", word.Occurrences, word.Word)
  }
}