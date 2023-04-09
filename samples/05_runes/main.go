package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	// Rune: an int32 that represents a Unicode code point for character
	// use runes if need to work with individual characters in a string
	rStr := "abcdefg"
	fmt.Println("Rune count -", utf8.RuneCountInString(rStr))

	fmt.Printf("%d - %#U - %v\n", 0, rStr[0], rStr[0])

	for index, val := range rStr {
		fmt.Printf("%d - %#U - %c\n", index, val, val)
	}
}
