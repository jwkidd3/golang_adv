package main

import "fmt"

func main() {
	// %d	:	Decimal Integer
	// %o	:	Base 8
	// %x	: 	Base 16 | Hexadecimal integer (lowercase letters)
	// %X	: 	Base 16 | Hexadecimal integer (uppercase letters)
	// %f	:	Float
	// %c	:	Character | rune | Unicode point
	// %U	:	Character's Unicode code point as a hexadecimal value (U+0041)
	// %#U	:	Character's Unicode code point as a hexadecimal value with that character (U+0041 'A')
	// %s	:	String
	// %q	:	Quoted string
	// %t	:	Boolean
	// %p	:	Pointer address
	// %v	:	Guesses based on data type
	// %T	:	Type of supplied value

	fmt.Printf("%s - %d - %o - %x - %f - %t\n", "hihi", 4, 5, 5, 3.14, true)

	char := 'A'
	fmt.Printf("Character: %c - %U - %#U\n", char, char, char)

	fmt.Printf("%9f\n", 3.14)
	fmt.Printf("%.f\n", 3.541592)
	fmt.Printf("%9.f\n", 3.141592)
	fmt.Printf("%.2f\n", 3.141592)

	sp1 := fmt.Sprintf("%.2f", 3.141592)
	fmt.Println(sp1)
}
