package words_test

import (
	"github.com/sethpollen/sbp_linux_config/sbpgo/games/words/coca"
	"testing"
)

func TestGetWordList(t *testing.T) {
	list := coca.GetWordList()

	// The expectation here is less than 5000 due to repeated words in the data
	// file.
	if len(list.Words) != 4351 {
		t.Errorf("Expected 4351 words; got %v", len(list.Words))
	}

	var totalOccurrences int64 = 0
	var lastOccurrences int64 = -1
	var used = make(map[string]bool)

	for _, word := range list.Words {
		totalOccurrences += word.Occurrences

		if lastOccurrences != -1 {
			if word.Occurrences > lastOccurrences {
				t.Error("WordList not sorted")
				return
			}
		}
		lastOccurrences = word.Occurrences

		if word.Word == "" {
			t.Error("Empty word")
			return
		}

		if used[word.Word] {
			t.Error("Duplicate word: ", word)
			return
		}
		used[word.Word] = true
	}

	if list.TotalOccurrences != totalOccurrences {
		t.Errorf("Expected %v total occurrences; got %v",
			totalOccurrences, list.TotalOccurrences)
	}
}
