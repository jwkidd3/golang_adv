/*
	A CLOSURE is a function value that references variables from outside its body.
	(
		Closure is a function that doesn't have to be associated with an identifier,
		but they are often associated with a variable.
		Closure can change values outside of the function.
	)

	The function may access and assign to the referenced variables;
	in this sense the function is "bound" to the variables.
*/

package main

import "fmt"

var pl = fmt.Println

func useFunc(f func(int, int) int, x, y int) {
	pl("Result :", f(x, y))
}

func sumVals(x, y int) int {
	return x + y
}

func main() {
	// Create a closure that sums values
	intSum := func(x, y int) int {
		return x + y
	}
	pl("4 + 2 =", intSum(4, 2))

	// Closures can change values outside the function
	val := 2
	changeVar := func() {
		val += 1
	}
	changeVar()
	pl("val =", val)

	// Pass function to another function
	useFunc(sumVals, 5, 6)
}
