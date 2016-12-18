package words

type Word struct {
	Word        string
	Occurrences int64
	// A single word may function as several parts of speech, so this string
	// may contain several characters. Each character refers to a different
	// part of speech. This string will be empty if the part of speech is
	// unknown. Currently, only the COCA corpus provides part of speech data.
	PartsOfSpeech string
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
