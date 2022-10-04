package misc

import (
	"regexp"
	"strings"
)

/*

  File:    soundex.go
 Author (Original:  Usman Mahmood
           Copyright (c) 2015 Usman Mahmood
           All rights reserved.
           MIT License
           https://github.com/umahmood/soundex

  Copyright (c) 2022. BSD 3-Clause License
	https://opensource.org/licenses/BSD-3-Clause

  The this permission notice shall be included in all copies
    or substantial portions of the Software.

*/
/*

Soundex is a go implementation of the Soundex algorithm.

"Soundex is a phonetic algorithm for indexing names by sound, as pronounced in English.
The goal is for homophones to be encoded to the same representation
so that they can be matched despite minor differences in spelling."
   - Wikipedia
Based on:
$ go get github.com/umahmood/soundex

*/

const EmptySoundex = "000"

// soundexMappings soundex consonant mappings
var soundexMappings = map[string][]string{
	"1": {"b", "f", "p", "v"},
	"2": {"c", "g", "j", "k", "q", "s", "x", "z"},
	"3": {"d", "t"},
	"4": {"l"},
	"5": {"m", "n"},
	"6": {"r"},
}

var soundexTable = make(map[string]string)

// transpose creates a lookup table of soundex mappings
func transpose() {
	if len(soundexTable) < 1 {
		soundexTable = make(map[string]string)
		for val, list := range soundexMappings {
			for _, char := range list {
				soundexTable[char] = val
			}
		}
	}
	return
}

var digit = regexp.MustCompile("[0-9]")
var ignore = regexp.MustCompile("[hw]*")
var vowel = regexp.MustCompile("[aeiouy]*")
var first = regexp.MustCompile("[1-6]")

// SoundexCode generates the soundex code for a given word
//goland:noinspection GoUnusedExportedFunction
func SoundexCode(word string) string {
	// remove numbers and spaces
	word = digit.ReplaceAllString(word, "")
	word = strings.Replace(word, " ", "", -1)
	if word == "" {
		return EmptySoundex
	}
	// transpose the mappings to create a lookup table for convenience
	transpose()
	// convert the word to lower case for consistent lookups
	word = strings.ToLower(word)
	// save the first letter of the word
	firstLetter := string(word[0])
	// remove all occurrences of 'h' and 'w' except first letter
	code := firstLetter + ignore.ReplaceAllString(word[1:], "")
	// replace all consonants (include the first letter) with digits based
	// on Soundex mapping table
	val := ""
	for idx, ch := range code {
		t := soundexTable[string(ch)]
		if t == "" {
			val += string(code[idx])
		} else {
			val += t
		}
	}
	code = val
	// replace all adjacent same digits with one digit
	n := strings.Split(code, "")
	a := 0
	b := 1
	for b < len(n) {
		if n[a] == n[b] {
			n = append(n[:a], n[a+1:]...)
		}
		a++
		b++
	}
	code = strings.Join(n, "")
	// remove all occurrences of a, e, i, o, u, y except first letter
	code = string(code[0]) + vowel.ReplaceAllString(code[1:], "")
	// if first symbol is a digit replace it with the saved first letter
	code = first.ReplaceAllString(string(code[0]), firstLetter) + code[1:]
	// append 3 zeros if result contains less than 3 digits
	if len(code) <= 3 {
		code += EmptySoundex
	}
	// retain the first 4 characters
	code = code[:4]
	// convert to upper case and return
	return strings.ToUpper(code)
}
