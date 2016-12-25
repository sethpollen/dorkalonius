// Library for de-inflecting words. For instance, we reduce verbs to their
// infinitive form and nouns to their singular form.

package wiktionary

import (
  "log"
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

func (self *InflectionMap) NumBaseWords() int {
	return len(self.BaseWords)
}

// Adds a baseWord, inflected pair to the map.
func (self *InflectionMap) Add(baseWord, inflected string) {
	self.BaseWords[baseWord] = true
	existingBaseWord, ok := self.InflectedToBase[inflected]
	if ok && existingBaseWord != baseWord {
    log.Fatalf("Multiple base words (%q, %q) for inflected form %q\n",
               existingBaseWord, baseWord, inflected)
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
