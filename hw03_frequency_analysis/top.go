package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

const OutputSliceLength = 10

type WordStat struct {
	Word  string
	Count int
}

var (
	reClearWord         = regexp.MustCompile(`[^\p{L}-]+`)
	reHasLetter         = regexp.MustCompile(`\p{L}`)
	AsteriskIsCompleted = true
)

// Top10 returns 10 most common words of the string.
func Top10(str string) []string {
	mappedWords := make(map[string]int)
	for _, word := range strings.Fields(str) { // count all words
		if AsteriskIsCompleted {
			word = strings.ToLower(reClearWord.ReplaceAllString(word, ""))
			if reHasLetter.MatchString(word) {
				mappedWords[word]++
			}
		} else {
			mappedWords[word]++
		}
	}
	ws := make([]WordStat, 0)
	for word, cnt := range mappedWords { // gather up the statistic
		ws = append(ws, WordStat{Word: word, Count: cnt})
	}
	sort.SliceStable(ws, func(i, j int) bool { // sorting lexicographically
		return ws[i].Word < ws[j].Word
	})
	sort.SliceStable(ws, func(i, j int) bool { // sorting by count
		return ws[i].Count > ws[j].Count
	})
	resultLen := min(OutputSliceLength, len(ws))
	result := make([]string, 0)
	for _, el := range ws[:resultLen] {
		result = append(result, el.Word)
	}
	return result
}

// min returns the entry in the list with the smallest numerical value.
func min(el int, elements ...int) int {
	min := el
	for _, el := range elements {
		if el < min {
			min = el
		}
	}
	return min
}
