// Library for de-inflecting words. For instance, we reduce verbs to their
// infinitive form and nouns to their singular form.

package wiktionary

// TODO: We may want to take part of speech into account when deciding which
// base form to choose among multiple base forms.

type InflectionMap struct {
	BaseWords       map[string]bool
	InflectedToBase map[string]string
}

func NewInflectionMap() *InflectionMap {
	return &InflectionMap{make(map[string]bool), make(map[string]string)}
}

func (self *InflectionMap) NumBaseWords() int {
	return len(self.BaseWords)
}

// Adds a baseWord, inflected pair to the map.
func (self *InflectionMap) Add(baseWord, inflected string) {
	self.BaseWords[baseWord] = true
	existingBaseWord, ok := self.InflectedToBase[inflected]
	if ok && len(existingBaseWord) <= len(baseWord) {
		// The existing base word is shorter, so just leave it.
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