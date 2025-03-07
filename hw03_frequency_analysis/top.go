package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

func Top10(str1 string) []string {
	var top1 []string
	var s []string = strings.Fields(str1) // слайс из исходной строки
	var s2 []int                          // слайс для подсчёта частоты слова
	var s3 []string                       // слайс уникальных значений
	for i, val1 := range s {
		index := 0
		check := true
		if i == 0 {
			s3 = append(s3, val1)
			for _, val2 := range s {
				if val2 == val1 {
					index++
				}
			}
			s2 = append(s2, index)
		}
		for _, valU := range s3 { // проверка на уникальность
			if val1 == valU {
				check = false
			}
		}
		if check { // Если уникально
			s3 = append(s3, val1)
			for _, val2 := range s {
				if val2 == val1 {
					index++
				}
			}
			s2 = append(s2, index)
		}
	}
	type WordAndIndex struct { // структура для сортировки с парами "слово"-"частота встреч"
		word  string
		count int
	}
	var worsAndIndexes []WordAndIndex // сортируем попарно по частоте встреси и в случае совпадения частоты встречи лексикографически
	for id := range s2 {
		worsAndIndexes = append(worsAndIndexes, WordAndIndex{word: s3[id], count: s2[id]})
	}
	sort.Slice(worsAndIndexes, func(i, j int) bool {
		if worsAndIndexes[i].count == worsAndIndexes[j].count {
			return worsAndIndexes[i].word < worsAndIndexes[j].word
		}
		return worsAndIndexes[i].count > worsAndIndexes[j].count
	})
	var k int
	if len(worsAndIndexes) < 10 {
		k = len(worsAndIndexes)
	} else {
		k = 10
	}
	for i := 0; i < k; i++ { // формируем выходной слайс длиной k (или 10 или равный количеству уникальных слов, в случае если их меньше 10)
		top1 = append(top1, worsAndIndexes[i].word)
	}
	return top1
}
