package main

import "fmt"

var pl = fmt.Println

func main() {
	// Conditional Operators:	< > <= >= == !=
	// Logical Operators:	&& || !

	age := 17
	if (age >= 1) && (age <= 18) {
		pl("Important Birthday")
	} else if (age == 21) || (age == 50) {
		pl("Important Birthday")
	} else if age >= 65 {
		pl("Important Birthday")
	} else {
		pl("Not Important Birthday")
	}

	pl("!true =", !true)
}
