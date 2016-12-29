package counter_test

import (
	"strings"
	"testing"
)
import . "github.com/sethpollen/dorkalonius/counter"

func TestEmpty(t *testing.T) {
	var counter int = 0
	ProcessWords(strings.NewReader(""), func(string) error {
		counter++
		return nil
	})
	if counter != 0 {
		t.Errorf("counter: %d", counter)
	}
}

func TestWords(t *testing.T) {
	var counter int = 0
	ProcessWords(strings.NewReader("Hey, :joe! 89"), func(word string) error {
		switch counter {
		case 0:
			if word != "hey" {
				t.Errorf("word: %q", word)
			}
			break
		case 1:
			if word != "joe" {
				t.Errorf("word: %q", word)
			}
			break
		default:
			t.Errorf("word: %q", word)
			break
		}
		counter++
		return nil
	})
	if counter != 2 {
		t.Errorf("counter: %d", counter)
	}
}
