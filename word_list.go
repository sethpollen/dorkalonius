package dorkalonius

type Word struct {
	Word        string
	Occurrences int64
	// Indicates whether this word is known to be used as an adjective. May be
	// false if we don't know the part of speech for this word. Currently, only
	// the COCA corpus provides part of speech data.
	Adjective   bool
}

type WordList struct {
	// Sorted by descending occurrence count.
	Words            []Word
	TotalOccurrences int64
}

// Constructs an empty WordList.
func NewWordList() *WordList {
	return &WordList{make([]Word, 0), 0}
}

func (self *WordList) AddWord(word Word) {
	self.Words = append(self.Words, word)
	self.TotalOccurrences += word.Occurrences
}

// Support for sorting WordList objects by descending occurrence count.
func (self *WordList) Len() int {
	return len(self.Words)
}
func (self *WordList) Swap(i, j int) {
	self.Words[i], self.Words[j] = self.Words[j], self.Words[i]
}
func (self *WordList) Less(i, j int) bool {
	return self.Words[i].Occurrences > self.Words[j].Occurrences
}
