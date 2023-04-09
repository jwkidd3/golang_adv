/*
	REGEX - REGULAR EXPRESSIONS
		a sequence of characters that forms a search pattern
		It is a powerful tool used for matching patterns in text
*/

package main

import (
	"fmt"
	"regexp"
)

var (
	pl              = fmt.Println
	isValidUsername = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
)

func main() {
	// ----- Way 1 -----
	rStr1 := "The apple tree is in the Apple-Company"
	// Find the word "apple" without space around
	match, _ := regexp.MatchString("(apple[^ ]?", rStr1)
	pl("Have only `apple` :", match)

	// ----- Way 2 -----
	rStr2 := "Cat rat mat fat pat"
	r, _ := regexp.Compile("([crmfp]at)")
	pl("MatchString :", r.MatchString(rStr2))
	pl("FindString :", r.FindString(rStr2))
	pl("Indexes :", r.FindAllStringIndex(rStr2, -1))
	pl("All string :", r.FindAllString(rStr2, -1))
	pl("First 2 strings :", r.FindAllString(rStr2, 2))
	pl("All submatch index :", r.FindAllStringSubmatchIndex(rStr2, -1))
	pl(r.ReplaceAllString(rStr2, "Dog"))

	// ----- Way 3 -----
	if !isValidUsername("some-username") {
		pl("Username must contain only letters, digits, or underscore")
	}
}
