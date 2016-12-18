// Go wrapper for the Wiktionary Lua scripts for English inflections. We
// could use something like
// https://github.com/Shopify/go-lua/blob/master/README.md to execute Lua
// within our Go program, but for the time being we just shell out to the
// Lua command.

package wiktionary

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

// TODO: add a test case which passes now, then add more as corner cases
// are discovered in the data dump

// This is safe to use from multiple goroutines.
type Inflector struct {
	// Directory where the Lua files may be found.
	luaDir string
}

const mainLuaScript = "en-headword.lua"

func checkForLuaScript(dir string) (bool, error) {
	cwd, err := os.Open(dir)
	if err != nil {
		return false, err
	}
	dirList, err := cwd.Readdir(-1)
	if err != nil {
		return false, err
	}
	for _, fileInfo := range dirList {
		if fileInfo.Name() == mainLuaScript {
			return true, nil
		}
	}
	return false, nil
}

func NewInflector() (*Inflector, error) {
	// We must search for the Lua files, as binaries and tests run in different
	// directories. In binaries, the Lua files may be found at
	// ./sbpgo/games/words/wiktionary. In unit tests, the Lua files are in the
	// current directory (.).
	for _, dir := range []string{".", "./sbpgo/games/words/wiktionary"} {
		found, err := checkForLuaScript(dir)
		if err != nil {
			return nil, err
		}
		if found {
			return &Inflector{dir}, nil
		}
	}
	return nil, errors.New("Could not find Lua scripts")
}

// Parts of speech which can be passed to ExpandInflections.
const (
	Noun = iota
	Verb
	Adjective
	Adverb
)

func PosName(pos int) string {
	return []string{"noun", "verb", "adjective", "adverb"}[pos]
}

// 'title' should be the base word of the Wiktionary page. 'args' should be
// the args passed to the en-verb or en-noun template. The return value
// contains the expanded list of inflections, not including the original
// 'title' word.
func (self *Inflector) ExpandInflections(
	pos int, title string, args []string) ([]string, error) {
  if title == "principl" {
    // This is a very weird entry in Wiktionary. Just ignore it.
    return nil, nil
  }
    
	var err error

	var posStr string
	switch pos {
	case Noun:
		posStr = "nouns"
	case Verb:
		posStr = "verbs"
	case Adjective:
		posStr = "adjectives"
	case Adverb:
		posStr = "adverbs"
	default:
		return nil, fmt.Errorf("Unsupported part of speech: %v", pos)
	}

	cmdArgs := []string{mainLuaScript, posStr, title}
	for _, arg := range args {
		// Sometimes wiktionary inserts a : at the beginning of an argument. I don't
		// know why, but it looks like we are supposed to drop it to make Lua happy.
		if strings.HasPrefix(arg, ":*") {
			arg = arg[1:]
		}
		cmdArgs = append(cmdArgs, arg)
	}

	cmd := exec.Command("lua", cmdArgs...)
	cmd.Dir = self.luaDir

	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		if stderr.Len() == 0 {
			return nil, err
		}
		return nil, errors.New(stderr.String())
	}

	results := make(map[string]bool)
	var debugLines []string

	// Split the stdout into lines.
	scanner := bufio.NewScanner(strings.NewReader(stdout.String()))
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("Could not parse Lua output line: \"%s\"",
				scanner.Text())
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		switch key {
		case "debug":
			debugLines = append(debugLines, value)
		default:
			for _, result := range processLuaResult(key, value) {
				results[result] = true
			}
		}
	}

	if len(debugLines) > 0 {
		log.Println("Lua debug logs:\n" + strings.Join(debugLines, "\n"))
	}

	// No need to include inflections which are the same as the base word.
	delete(results, title)

	var resultsList []string
	for result, _ := range results {
    if err = checkResult(result); err != nil {
      return nil, err
    }
		resultsList = append(resultsList, result)
	}
	return resultsList, nil
}

func processLuaResult(key, value string) []string {
	for _, prefix := range []string{"[[more]] ", "[[most]] "} {
		if strings.HasPrefix(value, prefix) {
			value = value[len(prefix):]
			break
		}
	}

	var results []string

	// Drop any multi-word forms, as we only process corpora one word at a
	// time.
	if strings.Index(value, " ") >= 0 {
		return results
	}

	results = append(results, value)

	// Include the plural form of the present participle, e.g. "wanderings".
	if key == "present-participle-form-of" &&
		strings.HasSuffix(value, "ing") {
		results = append(results, value+"s")
	}

	return results
}

func checkResult(result string) error {
  for r := range result {
    switch r {
      case '{', '}', '[', ']':
        return fmt.Errorf(
          "Got a result with MediaWiki template characters: %q", result)
    }
  }
  return nil
}
