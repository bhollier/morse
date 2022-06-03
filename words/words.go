package words

import (
	"bufio"
	_ "embed"
	"strings"
)

const numWords = 3000

//go:embed words.txt
var wordsFile []byte

var words []string

var runes = []rune("abcdefghijklmnopqrstuvwxyz")

var wordsWithRuneMap = make(map[rune][]string, len(runes))

func init() {
	words = make([]string, 0, numWords)
	scanner := bufio.NewScanner(strings.NewReader(string(wordsFile)))
	for scanner.Scan() {
		words = append(words, strings.ToLower(scanner.Text()))
	}

	for _, r := range runes {
		wordsWithRune := make([]string, 0)
		for _, word := range words {
			if strings.ContainsRune(word, r) {
				wordsWithRune = append(wordsWithRune, word)
			}
		}
		wordsWithRuneMap[r] = wordsWithRune
	}
}

// WithAnyRunes returns all words containing any of the given runes
func WithAnyRunes(rs []rune) (words []string) {
	wordSet := make(map[string]struct{})
	for _, r := range rs {
		for _, word := range wordsWithRuneMap[r] {
			wordSet[word] = struct{}{}
		}
	}
	words = make([]string, 0, len(wordSet))
	for word := range wordSet {
		words = append(words, word)
	}
	return
}

// WithOnlyRunes returns all the words that are made up of only the given runes
func WithOnlyRunes(rs []rune) []string {
	rSet := make(map[rune]struct{}, len(rs))
	for _, r := range rs {
		rSet[r] = struct{}{}
	}

	matchingWords := make([]string, 0)
	for _, word := range words {
		hasAllRunes := true
		for _, wr := range []rune(word) {
			_, ok := rSet[wr]
			if !ok {
				hasAllRunes = false
				break
			}
		}
		if hasAllRunes {
			matchingWords = append(matchingWords, word)
		}
	}
	return matchingWords
}
