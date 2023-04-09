/*	POINTERS
	A pointer holds the memory address of a value.

	The type *T is a pointer to a T value. Its zero value is nil.

	* is the dereferencing operator: gives the value of the pointer
	& is the address-of operator: gives the address of pointer

	pointer := new(datatype)
	pointer := &aVariable
*/

package main

import "fmt"

var pl = fmt.Println

func notChangeVal(f1 int) int {
	f1 += 1
	return f1
}

func changeVal(ptr *int) {
	*ptr = 12
}

// Not change the argument value - input array
func doubleArrVals(arr [4]int) {
	for index := 0; index < 4; index++ {
		arr[index] *= 2
	}
}

// Doing with input array pointer can change the array outside of the function
func doubleArrValsByPtr(arr *[4]int) {
	for index := 0; index < 4; index++ {
		arr[index] *= 2
	}
}

// Not change the argument value - input slice
func doubleVals(nums ...float64) {
	for _, val := range nums {
		val *= 2
	}
}

// Change the input slice
func doubleSliceVals(nums []float64) {
	for index := range nums {
		nums[index] *= 2
	}
}

func main() {
	f1 := 5
	pl("f1 before func notChangeVal :", f1)
	notChangeVal(f1)
	pl("f1 after func notChangeVal :", f1)

	f2 := 10
	// f2Ptr := &f2
	var f2Ptr *int = &f2
	pl("f2 address : ", f2Ptr)
	pl("f2 value : ", *f2Ptr)

	*f2Ptr = 11
	pl("f2 value : ", *f2Ptr)

	pl("f2 before func changeVal :", f2)
	changeVal(&f2)
	pl("f2 after func changeVal :", f2)

	// Pass array into function
	pArr := [4]int{1, 2, 3, 4}
	pl("pArr : ", pArr)
	doubleArrVals(pArr)
	pl("pArr after doubleArrVals :", pArr)
	doubleArrValsByPtr(&pArr) // Pass array pointer can change its values
	pl("pArr after doubleArrValsByPtr :", pArr)

	// Pass slice into function
	sl := []float64{11, 22, 33}
	pl("Slice :", sl)

	doubleVals(sl...)
	pl("Slice after doubleVals :", sl)
	doubleSliceVals(sl)
	pl("Slice after doubleSliceVals :", sl)
}
