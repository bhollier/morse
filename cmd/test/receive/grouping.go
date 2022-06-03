package receive

import (
	"fmt"
	"github.com/bhollier/morse/words"
	"math/rand"
	"strings"
	"unicode"
)

type Grouping rune

const (
	CharacterGrouping = Grouping('c')
	WordGrouping      = Grouping('w')
	SentenceGrouping  = Grouping('s')
)

func ParseGrouping(s string) (Grouping, error) {
	if s == "" {
		return 0, fmt.Errorf("missing group")
	}

	switch Grouping(unicode.ToLower([]rune(s)[0])) {
	case CharacterGrouping:
		return CharacterGrouping, nil
	case WordGrouping:
		return WordGrouping, nil
	case SentenceGrouping:
		return SentenceGrouping, nil
	default:
		return 0, fmt.Errorf("unknown group %s", s)
	}
}

func randomRune(r *rand.Rand) rune {
	runes := []rune(*characters)
	return runes[r.Intn(len(runes))]
}

func generateRandomWord(r *rand.Rand) string {
	wordsWithRunes := words.WithOnlyRunes([]rune(*characters))
	// Filter out words that are too long
	i := 0
	for _, word := range wordsWithRunes {
		if uint(len(word)) <= *maxWordLength {
			wordsWithRunes[i] = word
			i++
		}
	}
	wordsWithRunes = wordsWithRunes[:i]
	// Return a random word
	return wordsWithRunes[r.Intn(len(wordsWithRunes))]
}

func generateRandomSentence(r *rand.Rand) string {
	sentenceLen := r.Intn(int(*maxSentenceLength-*minSentenceLength)) + int(*minSentenceLength)
	sb := strings.Builder{}
	for i := 0; i < sentenceLen; i++ {
		sb.WriteString(generateRandomWord(r))
		sb.WriteRune(' ')
	}
	return sb.String()
}

func (g Grouping) GenerateString(r *rand.Rand) string {
	switch g {
	case CharacterGrouping:
		return string([]rune{randomRune(r)})
	case WordGrouping:
		return generateRandomWord(r)
	case SentenceGrouping:
		return generateRandomSentence(r)
	default:
		panic(fmt.Errorf("unknown group mode %c", g))
	}
}
