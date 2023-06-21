package hw03frequencyanalysis

import (
	"container/heap"
	"sort"
	"strings"
)

const countOfTopElements = 10

func Top10(input string) []string {
	inputSlice := strings.Fields(input)
	if len(inputSlice) < countOfTopElements {
		return []string{}
	}

	wordFrequencyMap := make(map[string]int)
	for _, word := range inputSlice {
		wordFrequencyMap[word]++
	}

	topElements := newMinHeap()

	for word, frequency := range wordFrequencyMap {
		heap.Push(topElements, Pair{word, frequency})
		if topElements.Len() > countOfTopElements {
			heap.Pop(topElements)
		}
	}

	topPairs := make([]Pair, 0)
	for !topElements.Empty() {
		pair := topElements.Top()
		heap.Pop(topElements)
		topPairs = append(topPairs, pair)
	}

	sort.Slice(topPairs, func(i, j int) bool {
		if topPairs[i].frequency != topPairs[j].frequency {
			return topPairs[i].frequency > topPairs[j].frequency
		}
		return topPairs[i].word < topPairs[j].word
	})

	topWords := make([]string, 0, len(topPairs))
	for _, pair := range topPairs {
		topWords = append(topWords, pair.word)
	}

	return topWords
}
