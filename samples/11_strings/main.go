package main

import (
	"fmt"
	"strings"
)

var pl = fmt.Println
var pf = fmt.Printf

func main() {
	// Strings are immutable | array of bytes | array of runes
	s1 := "abcd"
	rArr := []rune(s1)
	pl("Rune array:")
	for _, val := range rArr {
		pf("\t %c - %d\n", val, val)
	}

	bArr := []byte{'a', 'b', 'c', 'd'}
	bStr := string(bArr[:])
	pl("I am:", bArr)
	pl("I am:", bStr)

	sV1 := "A word"
	replacer := strings.NewReplacer("A", "Another")
	sV2 := replacer.Replace(sV1)
	pl("sV2 - ", sV2)

	pl("Length -", len(sV2))
	pl("Contains Another -", strings.Contains(sV2, "Another"))
	pl("o index -", strings.Index(sV2, "o"))
	pl("Replace o -> 0 -", strings.Replace(sV2, "o", "0", -1)) // -1 : no limit on the number of replacements

	sV3 := "\nSome words\n" // \t \" \\
	sV4 := strings.TrimSpace(sV3)
	pf("Trim '%v' - %v\n", sV3, sV4)

	pl("Split -", strings.Split("a-b-c-d", "-"))
	pl("Lower -", strings.ToLower(sV4))
	pl("Upper -", strings.ToUpper(sV4))
	pl("Prefix -", strings.HasPrefix("awesome", "aw"))
	pl("Suffix -", strings.HasSuffix("awesome", "some"))
}
