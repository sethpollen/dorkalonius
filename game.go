// Basic game logic.

package dorkalonius

type Game struct {
	TargetWord     string
	AvailableWords []string
}

const (
	numAvailableWords = 35

	// Manual tuning parameters. We use a high bias for the target word
	// in order to get something interesting. We use a much smaller bias
	// for the available words, since we want them to mostly reflect a
	// typical selection of words.
	// TODO: consider adjusting the targetWordBias
	targetWordBias    float64 = 3e-4
	availableWordBias float64 = 3e-6
)

func NewGame(wordSet WordSet) (*Game, error) {
	adjectives := GetCocaAdjectives()
	adjective := adjectives.Sample(
		1, int64(targetWordBias*float64(adjectives.Size()))).GetWords()

	words := wordSet.Sample(numAvailableWords,
		int64(availableWordBias*float64(wordSet.Size())))
	wordsSlice := words.GetWords()
	bareWords := make([]string, len(wordsSlice))
	for i := range wordsSlice {
		bareWords[i] = wordsSlice[i].Word
	}

	return &Game{adjective[0].Word, bareWords}, nil
}
