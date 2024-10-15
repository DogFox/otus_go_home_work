package hw03frequencyanalysis

import (
	"sort"
	"strings"
	"unicode"
)

type WordCount struct {
	word  string
	count int
}

func splitBySymbolsWithoutDash(c rune) bool {
	return !unicode.IsLetter(c) && !unicode.IsNumber(c) && c != '-'
}

func prepareMap(str string) map[string]int {
	words := strings.FieldsFunc(str, splitBySymbolsWithoutDash)
	freqMap := make(map[string]int)
	// fmt.Println(words)

	for _, v := range words {
		lowerCaseWord := strings.ToLower(v)
		freqMap[lowerCaseWord]++
	}
	return freqMap
}

func Top10(str string) []string {
	words := prepareMap(str)

	slice := make([]WordCount, len(words))
	i := 0
	for k, v := range words {
		slice[i].word = k
		slice[i].count = v
		i++
	}

	sort.Slice(slice, func(i, j int) bool {
		// словом не считаем -> на выселки
		if slice[i].word == string('-') {
			return false
		}
		// сортируем по колву
		if slice[i].count != slice[j].count {
			return slice[i].count > slice[j].count
		}
		// сортируем по словарю
		return slice[i].word < slice[j].word
	})

	if len(slice) > 10 {
		slice = slice[0:10]
	}
	result := make([]string, len(slice))

	for i, v := range slice {
		result[i] = v.word
	}

	// fmt.Println(result)
	return result
}
