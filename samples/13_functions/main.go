package main

import "fmt"

var pl = fmt.Println
var pf = fmt.Printf

// === FUNCTION ===
// func funcName(parameters) returnTypes {BODY}

// Named function
func sayHello() {
	pl("Hello")
}

// func getSum(x int, y int) int {
func getSum(x, y int) int {
	return x + y
}

// Multiple parameters and multiple return values
func getQuotient(x float64, y float64) (float64, error) {
	if y == 0 {
		return 0, fmt.Errorf("Can not divide by zero")
	}

	return x / y, nil
}

// Both the input parameters and return values can be named
func add(x, y int) (sum int, diff int) {
	sum = x + y
	diff = x - y

	// Function returns named return values without specifying
	return
}

// Varadic Functions
// a function receives an unknown the number of parameters
func getSum2(nums ...int) int {
	sum := 0
	for _, num := range nums {
		sum += num
	}
	return sum
}

// Doing with input slice does not change the slice outside of the function
func getSliceSum(sl []int) int {
	sum := 0
	for _, val := range sl {
		sum += val
	}
	return sum
}

func main() {
	sayHello()
	pl("getSum(4, 5) :", getSum(4, 5))

	result, err := getQuotient(2, 0)
	pf("getQuotient(2, 0) : %v %v\n", result, err)

	sum, diff := add(1, 2)
	pf("add(1, 2) : %d %d\n", sum, diff)
	pl("getSum2(1, 3, 5, 7) :", getSum2(1, 3, 5, 7))

	// Pass slice into functions
	sl := []int{1, 2, 3, 4}
	pl("getSliceSum :", getSliceSum(sl))
	pl("getSum2(sl...) :", getSum2(sl...))
}
