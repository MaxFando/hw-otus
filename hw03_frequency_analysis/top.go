package hw03frequencyanalysis

import (
	"container/heap"
	"strings"
)

const countOfTopElements = 10

func Top10(input string) []string {
	inputSlice := strings.Fields(input)

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

	topWords := make([]string, topElements.Len())
	i := topElements.Len() - 1
	for !topElements.Empty() {
		pair := heap.Pop(topElements)
		topWords[i] = pair.(Pair).word
		i--
	}

	return topWords
}
