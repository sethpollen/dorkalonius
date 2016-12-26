// Library for de-inflecting words. For instance, we reduce verbs to their
// infinitive form and nouns to their singular form.

package wiktionary

import (
	"bufio"
	"compress/bzip2"
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"io"
  "log"
	"os"
	"strings"
)

type InflectionMap struct {
	BaseWords       map[string]bool
	InflectedToBase map[string]string

	// Manually selected overrides to use when populating InflectedToBase; this
	// works around inflected forms which map to multiple base forms.
	PreferredInflectedToBase map[string]string
}

// Returns a new InflectionMap, initialized with 'data'.
func NewInflectionMap(data []Inflection,
	preferences map[string]string) *InflectionMap {
	m := &InflectionMap{make(map[string]bool),
		make(map[string]string),
		preferences}
	for _, i := range data {
		if i.Pos != "noun" && i.Pos != "verb" {
			// We don't deconjugate adjectives or adverbs. Doing so would map
			// comparatives and superlatives back to their base form; we have
			// decided not to do that.
			continue
		}
		m.Add(i.BaseWord, i.InflectedForms)
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

	preferences, err := loadPreferences()
	if err != nil {
		return nil, err
	}

	return NewInflectionMap(parsed.Inflections, preferences), nil
}

func (self *InflectionMap) NumBaseWords() int {
	return len(self.BaseWords)
}

func (self *InflectionMap) Add(baseWord string, inflectedForms []string) {
	self.BaseWords[baseWord] = true

	for _, inflected := range inflectedForms {
		if inflected == "-" || inflected == "?" {
			continue
		}

		// First check if we have a preference for this inflected form. If we do, then
		// always use that.
		preferredBaseWord, ok := self.PreferredInflectedToBase[inflected]
		if ok {
			self.InflectedToBase[inflected] = preferredBaseWord
			continue
		}

		existingBaseWord, ok := self.InflectedToBase[inflected]
		if !ok {
			self.InflectedToBase[inflected] = baseWord
			continue
		}
		if existingBaseWord == baseWord {
			continue
		}
		// We have conflicting base words for this inflected form.

		// Prefer to reduce -ings to -ing and not all the way down to the infinitive
		// form of the verb. Thus, "bearings" becomes "bearing" and not "bear".
		if strings.HasSuffix(inflected, "ings") {
			singular := inflected[0 : len(inflected)-1]
			if baseWord == singular {
				self.InflectedToBase[inflected] = baseWord
				continue
			}
			if existingBaseWord == singular {
				continue
			}
		}

		// TODO:
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Fprintf(os.Stderr,
				"(%d) Inflected %q maps to bases (%q, %q)? --> ",
        len(self.BaseWords), inflected, existingBaseWord, baseWord)
			chosenBase, _ := reader.ReadString('\n')
			chosenBase = strings.TrimSpace(chosenBase)
			if chosenBase != existingBaseWord && chosenBase != baseWord {
				fmt.Fprintf(os.Stderr, "Bad choice; try again\n")
				continue
			}

			self.InflectedToBase[inflected] = chosenBase
			self.PreferredInflectedToBase[inflected] = chosenBase
			break
		}

		fmt.Println("")
		for k, v := range self.PreferredInflectedToBase {
			fmt.Printf("%q,%q\n", k, v)
		}
		fmt.Println("")
	}
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

func loadPreferences() (map[string]string, error) {
	csvReader := csv.NewReader(Get_preferences_csv("preferences.csv"))
	result := make(map[string]string)
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return result, err
		}
		if len(record) != 2 {
			return result, fmt.Errorf(
				"Record has wrong number of cells: %d", len(record))
		}
		_, ok := result[record[0]]
		if ok {
      log.Fatalf("loadPreferences found duplicate entry for %q\n", record[0])
    }
		result[record[0]] = record[1]
	}
	return result, nil
}
