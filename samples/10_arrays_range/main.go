package main

import "fmt"

var pl = fmt.Println
var pf = fmt.Printf

func main() {
	// === Array ===
	// Array with fixed size (default)
	var arr1 [4]int
	arr1[0] = 1
	pl("arr1[0] -", arr1[0])
	pf("arr1 - %v - len(%d) - cap(%d)\n", arr1, len(arr1), cap(arr1))

	// Array literal can be specified like so
	arr2 := [5]int{1, 2, 3, 4, 5}
	pf("arr2 - %v - len(%d) - cap(%d)\n", arr2, len(arr2), cap(arr2))

	// The compiler count the array elements for you
	arr3 := [...]int{4, 5, 6}
	pf("arr3 - %v - len(%d) - cap(%d)\n", arr3, len(arr3), cap(arr3))

	for index := 0; index < len(arr3); index++ {
		pl(arr3[index])
	}

	for index, value := range arr3 {
		pf("%d : %d\n", index, value)
	}

	// 2D array
	arr4 := [2][2]int{
		{1, 2},
		{4, 5},
	}
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			pl(arr4[i][j])
		}
	}
}
