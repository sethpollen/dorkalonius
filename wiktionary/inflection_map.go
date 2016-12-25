// Library for de-inflecting words. For instance, we reduce verbs to their
// infinitive form and nouns to their singular form.

package wiktionary

import (
  "compress/bzip2"
  "encoding/xml"
  "log"
  "os"
  "strings"
)

type InflectionMap struct {
	BaseWords       map[string]bool
	InflectedToBase map[string]string
}

// Returns a new InflectionMap, initialized with 'data'.
func NewInflectionMap(data []Inflection) *InflectionMap {
	m := &InflectionMap{make(map[string]bool), make(map[string]string)}
	for _, i := range data {
    for _, inflectedForm := range i.InflectedForms {
      m.Add(i.BaseWord, inflectedForm)
    }
  }
  return m
}

func InflectionMapFromBzippedXml(filename string) (*InflectionMap, error) {
  var err error
  
  file, err := os.Open(filename)
  if err != nil {
    return nil, err
  }
  decoder := xml.NewDecoder(bzip2.NewReader(file))
  
  var parsed Inflections
  if err = decoder.Decode(&parsed); err != nil {
    return nil, err
  }
  
  return NewInflectionMap(parsed.Inflections), nil
}

func (self *InflectionMap) NumBaseWords() int {
	return len(self.BaseWords)
}

// Adds a baseWord, inflected pair to the map.
func (self *InflectionMap) Add(baseWord, inflected string) {
  if inflected == "-" {
    return
  }
  
	self.BaseWords[baseWord] = true
	existingBaseWord, ok := self.InflectedToBase[inflected]
	
	if ok && existingBaseWord != baseWord {
    // Prefer to reduce -ings to -ing and not all the way down to the infinitive
    // form of the verb. Thus, "bearings" becomes "bearing" and not "bear".
    if strings.HasSuffix(inflected, "ings") {
      if strings.HasSuffix(baseWord, "ing") {
        self.InflectedToBase[inflected] = baseWord
      }
      return
    }
    // TODO:
    log.Printf("Inflected %q maps to bases (%q, %q)\n",
               inflected, existingBaseWord, baseWord)
    return
	}
	self.InflectedToBase[inflected] = baseWord
}

// Gets the base word for the given inflected form.
func (self *InflectionMap) GetBaseWord(inflected string) string {
	// If the inflectedForm is itself a base word, do nothing.
	if self.BaseWords[inflected] {
		return inflected
	}
	baseWord, ok := self.InflectedToBase[inflected]
	if ok {
		return baseWord
	}
	// We don't have any mapping for this word, so just pass it through.
	return inflected
}

type preferenceSorter struct {
  data *[]string
}

func (self preferenceSorter) Len() int {
  return len(*self.data)
}

func (self preferenceSorter) Less(i, j int) bool {
  return len((*self.data)[i]) < len((*self.data)[j])
}

func (self preferenceSorter) Swap(i, j int) {
  temp := (*self.data)[i]
  (*self.data)[i] = (*self.data)[j]
  (*self.data)[j] = temp
}
