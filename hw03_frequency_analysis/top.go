package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(str1 string) []string {
	if len(str1) == 0 {
		return nil
	}
	var top1 []string
	s := strings.Fields(str1)            // слайс из исходной строки
	WordAndIndex := make(map[string]int) // мапа с уникальными словами и их количеством
	for _, word := range s {
		WordAndIndex[word]++
	}
	var s3 []string // слайс с уникальными словами
	for key := range WordAndIndex {
		s3 = append(s3, key)
	}
	sort.Slice(s3, func(i, j int) bool {
		if WordAndIndex[s3[i]] != WordAndIndex[s3[j]] { //сортировка по количеству совпадений
			return WordAndIndex[s3[i]] > WordAndIndex[s3[j]]
		}
		return s3[i] < s3[j] // сортировка лексикографическая
	})
	var k int
	if len(s3) < 10 {
		k = len(s3)
	} else {
		k = 10
	}
	for i := 0; i < k; i++ { // формируем выходной слайс длиной k
		top1 = append(top1, s3[i])
	}
	return top1
}
