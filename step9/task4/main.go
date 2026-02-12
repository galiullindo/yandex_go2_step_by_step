package main

import (
	"slices"
	"strings"
)

func AreAnagrams(word string, wordToo string) bool {
	word = strings.ToLower(word)
	wordToo = strings.ToLower(wordToo)

	if len(word) != len(wordToo) {
		return false
	}

	rWord := []rune(word)
	rWordToo := []rune(wordToo)

	slices.Sort(rWord)
	slices.Sort(rWordToo)

	return slices.Equal(rWord, rWordToo)
}
