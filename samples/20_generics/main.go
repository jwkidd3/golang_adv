/*
	GENERICS
	
	Go functions can be written to work on multiple types using type parameters.
	The type parameters of a function appear between brackets, before the function's arguments.
*/
package main

import "fmt"

var pl = fmt.Println

type MyConstraint interface {
	int | float64
}

func sum[T MyConstraint](x T, y T) T {
	return x + y
}

// comparable is a useful constraint that makes it possible to use the == and != operators on values of the type
func index[T comparable](sl []T, x T) int {
	for index, value := range sl {
		if value == x {
			return index
		}
	}

	return -1
}

func main() {
	pl("4 + 2 =", sum(4, 2))
	pl("4.6 + 2.5 =", sum(4.6, 2.5))

	sl := []int{1, 2, 3, 4}
	pl("index of 2 :", index(sl, 2))
}
