package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

var NotWordCharsRegexp = regexp.MustCompile(`[^\wа-яА-Я-]+`)

const WordsLimit = 10

type WordFrequency struct {
	word string
	freq int
}

func Top10(text string) []string {
	// удаление спец.символов и приведение к нижнему регистру
	text = NotWordCharsRegexp.ReplaceAllString(text, " ")
	text = strings.ToLower(text)

	// разделение на слова по 1 или более пробелам
	words := strings.Fields(text)

	// подсчёт частоты слов
	wordsCount := make(map[string]int)
	for _, word := range words {
		if word != "-" {
			wordsCount[word]++
		}
	}

	// стуктура для сортировки слов по частоте
	wordsFrequency := make([]WordFrequency, 0, len(wordsCount))
	for word, freq := range wordsCount {
		wordsFrequency = append(wordsFrequency, WordFrequency{word: word, freq: freq})
	}

	// лексикографическая сортировка
	sort.Slice(wordsFrequency, func(i, j int) bool {
		if wordsFrequency[i].freq == wordsFrequency[j].freq {
			return wordsFrequency[i].word < wordsFrequency[j].word
		}
		return wordsFrequency[i].freq > wordsFrequency[j].freq
	})

	// выбор топ-10
	top10 := make([]string, 0, WordsLimit)
	for i := 0; i < WordsLimit && i < len(wordsFrequency); i++ {
		top10 = append(top10, wordsFrequency[i].word)
	}

	return top10
}
