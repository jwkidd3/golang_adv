package main

import "fmt"

var pl = fmt.Println
var pf = fmt.Printf

func main() {
	// === Slice ===
	// Array with dynamic size

	// var name []datatype
	// var name []datatype = []datatype{v1, v2, ...}
	// name := []datatype{v1, v2, ...}
	// name := make([]datatype, len, cap)
	sl1 := make([]string, 6)
	sl1[0] = "I"
	sl1[1] = "am"
	sl1[2] = "iron"
	sl1[3] = "man"
	pf("Slice - %v - len(%d) - cap(%d)\n", sl1, len(sl1), cap(sl1))

	for index := 0; index < len(sl1); index++ {
		pl(sl1[index])
	}

	for _, val := range sl1 {
		pl(val)
	}

	// Slices are like references to arrays
	// A slice does not store any data, it just describes a section of an underlying array.
	// Changing the elements of a slice modifies the corresponding elements of its underlying array.
	// Other slices that share the same underlying array will see those changes.
	iArr := [5]int{1, 2, 3, 4, 5}
	sl2 := iArr[0:2]
	sl3 := iArr[1:2]
	pl("First 3 :", iArr[:3])
	pl("Last 3 :", iArr[2:])

	iArr[1] = 10
	pf("sl3 - %v - len(%d) - cap(%d) - address(%p)\n", sl3, len(sl3), cap(sl3), sl3)
	sl2[0] = 8
	pf("sl2 - %v - len(%d) - cap(%d) - address(%p)\n", sl2, len(sl2), cap(sl2), sl2)
	pf("iArr - %v - len(%d) - cap(%d) - address(%p)\n", iArr, len(iArr), cap(iArr), &iArr)

	sl3 = append(sl3, 20)
	pf("sl2 - %v - len(%d) - cap(%d) - address(%p)\n", sl2, len(sl2), cap(sl2), sl2)
	pf("sl3 - %v - len(%d) - cap(%d) - address(%p)\n", sl3, len(sl3), cap(sl3), sl3)
	pf("iArr - %v - len(%d) - cap(%d) - address(%p)\n", iArr, len(iArr), cap(iArr), &iArr)

	// Empty slice
	sl4 := make([]string, 5)
	pl("sl4 -", sl4)

	// nil value
	pl("sl4[0] -", sl4[0])
}
