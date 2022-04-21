package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

const topCount = 10

type WordCount struct {
	Word  string
	Count int
}

func Top10(str string) []string {
	if len(str) == 0 {
		return make([]string, 0)
	}

	wordSlice := strings.Fields(str)
	wordsCount := make(map[string]int)

	for _, word := range wordSlice {
		wordsCount[word]++
	}

	sortedWords := make([]WordCount, 0)
	for word, count := range wordsCount {
		sortedWords = append(sortedWords, WordCount{word, count})
	}

	sort.Slice(sortedWords, func(i, j int) bool {
		if sortedWords[i].Count == sortedWords[j].Count {
			return sortedWords[i].Word < sortedWords[j].Word
		}
		return sortedWords[i].Count > sortedWords[j].Count
	})

	res := make([]string, len(sortedWords))
	for i, word := range sortedWords {
		res[i] = word.Word
	}
	if len(res) >= topCount {
		res = res[0:topCount]
	}

	return res
}
