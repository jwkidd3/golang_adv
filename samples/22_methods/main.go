/*
	A method is a function with a special argument called a receiver.
	For a given type, the name of each method must be unique,

		and neither methods nor functions support overloading.

	type SomeStruct struct {}
	func (s *SomeStruct) doSomething() {}
*/
package main

import "fmt"

type Rectangle struct {
	length, height float64
}

// Method “ with receiver `r`
// Method is a function that is a part of the struct
func (r Rectangle) Area() float64 {
	return r.length * r.height
}

func main() {
	rect := Rectangle{5, 6}
	fmt.Println("Rect area :", rect.Area())
}
