package main

import "fmt"

// Create a new type with a underlying type of int
type age int

func main() {
	// Create a variable of type "age"
	var myAge age = 26

	fmt.Println("My age is", myAge)
}
