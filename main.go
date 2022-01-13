package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

// A container for a word and its score
type scoredWord struct {
	word  string
	score float64
}

//// quicksort zone: Copy pasted from this gist https://gist.github.com/imwally/58d6bb9bf9da098064054f73a19cdca1 ////
func partition(a []scoredWord, lo, hi int) int {
	p := a[hi].score
	for j := lo; j < hi; j++ {
		if a[j].score > p {
			a[j], a[lo] = a[lo], a[j]
			lo++
		}
	}

	a[lo], a[hi] = a[hi], a[lo]
	return lo
}

func quickSort(a []scoredWord, lo, hi int) {
	if lo > hi {
		return
	}

	p := partition(a, lo, hi)
	quickSort(a, lo, p-1)
	quickSort(a, p+1, hi)
}

//// end quicksort zone ////

// simple function to find max value in an array of 26 floats
func findMax(input [26]float64) (int, float64) {
	max := float64(0)
	maxIndex := 0
	for i, data := range input {
		if data > max {
			max = data
			maxIndex = i
		}
	}
	return maxIndex, max
}

// Finds a word's score, as defined by the sum of each letter's probability in its position
func wordScore(word string, wordProbabilities [5][26]float64) float64 {
	totalScore := 0.0
	for i, letter := range word {
		totalScore += wordProbabilities[i][int(letter-97)]
	}
	return totalScore
}

func main() {
	// List of words, source: http://aspell.net/
	dat, err := os.ReadFile("en_US.dic")
	if err != nil {
		log.Fatal(err)
	}

	// Get a clean slice of all five-letter words with no accents, that are not proper nouns (no capitals)
	allWordsRaw := strings.Split(string(dat), "\n")
	words := make([]string, 0)

	r, _ := regexp.Compile("^[a-z]{5}$")

	letterCounts := [5][26]float64{}

	for _, rawWord := range allWordsRaw {
		word := strings.Split(rawWord, "/")[0]
		if r.MatchString(word) {
			words = append(words, word)

			for i, letter := range word {
				letterCounts[i][letter-97] += 1
			}
		}
	}

	// Get probability of each letter at each position
	totalWords := float64(len(words))
	letterProabilities := [5][26]float64{}
	for i, position := range letterCounts {
		for j, count := range position {
			letterProabilities[i][j] = count / totalWords
		}
	}

	fmt.Printf("There are %v words\n\n", len(words))

	// Find most frequent letter at each position
	charToIntOffset := 97
	for i, position := range letterProabilities {
		max, score := findMax(position)
		fmt.Printf("In position %d, most frequent is: %c (%v%%)\n", i+1, byte(max+charToIntOffset), int(score*100))
	}

	// Find most highest-scoring words, as well as highest-scoring with no y in the final position
	sortedWords := []scoredWord{}
	sortedWordsNoY := []scoredWord{}
	for _, word := range words {
		score := wordScore(word, letterProabilities)
		sortedWords = append(sortedWords, scoredWord{word, score})
		if word[len(word)-1] != 'y' {
			sortedWordsNoY = append(sortedWordsNoY, scoredWord{word, score})
		}
	}
	quickSort(sortedWords, 0, len(sortedWords)-1)
	quickSort(sortedWordsNoY, 0, len(sortedWordsNoY)-1)

	fmt.Println("\nTop ten words (including ending in y):")
	for i, word := range sortedWords[0:10] {
		fmt.Printf("%v. %v\n", i+1, word.word)
	}

	fmt.Println("\nTop ten words (excluding ending in y):")
	for i, word := range sortedWordsNoY[0:10] {
		fmt.Printf("%v. %v\n", i+1, word.word)
	}
}
