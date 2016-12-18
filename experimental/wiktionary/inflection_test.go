package wiktionary_test

import (
	"sort"
	"strings"
	"testing"
)
import . "github.com/sethpollen/sbp_linux_config/sbpgo/games/words/wiktionary"

type Case struct {
	Pos   int
	Title string
	Args  []string
	// Expected results.
	Inflections []string
}

func MakeCases() []Case {
	return []Case{
		Case{Verb, "pound",
			[]string{},
			[]string{"pounds", "pounding", "poundings", "pounded"}},
		Case{Verb, "dictionary",
			[]string{"dictionar", "ies"},
			[]string{"dictionaries", "dictionarying", "dictionaryings",
				"dictionaried"}},
		Case{Noun, "dictionary",
			[]string{"dictionaries"},
			[]string{"dictionaries"}},
		Case{Verb, "free",
			[]string{"d"},
			[]string{"frees", "freeing", "freeings", "freed"}},
		Case{Noun, "free",
			[]string{},
			[]string{"frees"}},
		Case{Adjective, "free",
			[]string{"er"},
			[]string{"freer", "freest"}},
		Case{Adjective, "institutional",
			[]string{},
			[]string{}},
		Case{Adverb, "free",
			[]string{},
			[]string{}},
		Case{Verb, "cat",
			[]string{"catt"},
			[]string{"cats", "catting", "cattings", "catted"}},
		Case{Verb, "crow",
			[]string{"crows", "crowing", "crowed", "past2=crew", "past2_qual=UK",
				"crowed"},
			[]string{"crows", "crowing", "crowings", "crew", "crowed"}},
		Case{Verb, "carry",
			[]string{"ies"},
			[]string{"carries", "carrying", "carryings", "carried"}},
		Case{Noun, "*nix",
			[]string{":*nixes", ":*nices"},
			[]string{"*nixes", "*nices"}},
	}
}

func TestInflector(t *testing.T) {
	inflector, err := NewInflector()
	if err != nil {
		t.Fatal(err)
	}

	for _, c := range MakeCases() {
		actual, err := inflector.ExpandInflections(c.Pos, c.Title, c.Args)
		if err != nil {
			t.Fatalf("Got error: %v\n%v", err, c)
		}
		if len(actual) != len(c.Inflections) {
			t.Fatalf("Unexpected result: %v\n%v", strings.Join(actual, ", "), c)
		}

		sort.Strings(actual)
		sort.Strings(c.Inflections)

		for i := range actual {
			if actual[i] != c.Inflections[i] {
				t.Fatalf("Unexpected result: %v\n%v", strings.Join(actual, ", "), c)
				break
			}
		}
	}
}
