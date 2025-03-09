package hw03frequencyanalysis

import (
	"fmt"
	"sort"
	"strings"
)

func Top10(str1 string) []string {
	var top1 []string
	s := strings.Fields(str1) // слайс из исходной строки
	var s2 []int              // слайс для подсчёта частоты слова
	var s3 []string           // слайс уникальных значений
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
	WordAndIndex := make(map[string]int)
	for id := range s2 {
		WordAndIndex[s3[id]] = s2[id]
	}
	sort.Slice(s3, func(i, j int) bool {
		if WordAndIndex[s3[i]] != WordAndIndex[s3[j]] {
			return WordAndIndex[s3[i]] > WordAndIndex[s3[j]]
		}
		return s3[i] < s3[j]
	})
	fmt.Println(s3)
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
