/*
	RECURSION
		function calls itself
		There must be a condition to end the calling repetitively of that function
*/

package main

import "fmt"

func factorial(num int) int {
	if num == 0 {
		return 1
	}
	return num * factorial(num-1)
}

func main() {
	fmt.Println("3! =", factorial(3))
	// 1st: result = 3 * factorial(2) = | 3 * 2 △ = 6
	// 2nd:	result = 2 * factorial(1) = | 2 * 1 | = 2
	// 3rd: result = 1 * factorial(0) = ▽ 1 * 1 | = 1
}
