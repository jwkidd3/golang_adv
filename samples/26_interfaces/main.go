/*
	INTERFACE

	A contract that if anything inherits it must implement some defined methods.

	Abstract (in OOP):
		In Go, there is no keyword or concept of abstract classes as found in some other object-oriented languages such as Java
		However, we can achieve similar functionality using interfaces in Go.

	Polymorphism (in OOP):
		Go supports polymorphism through the use of interfaces
		Allows different objects to perform the same function in different ways.
*/

package main

import "fmt"

type Shape interface {
	Area() float64
}

type Rectangle struct {
	width  float64
	height float64
}

type Circle struct {
	radius float64
}

// Rectangle implements Shape's Area() method implicitly
func (r Rectangle) Area() float64 {
	return r.width * r.height
}

// Circle also implements Shape's Area() method with the different way - Polymorphism
func (c Circle) Area() float64 {
	return 3.14 * c.radius * c.radius
}

func PrintArea(s Shape) {
	fmt.Printf("Area of %T : %v\n", s, s.Area())
}

func (r *Rectangle) Perimeter() float64 {
	return 2 * (r.width + r.height)
}

func main() {
	r := Rectangle{width: 5, height: 6}
	c := Circle{radius: 28}
	PrintArea(r)
	PrintArea(c)

	var s1 Shape
	s1 = Rectangle{width: 4, height: 5}
	PrintArea(s1)

	// Cast the interface to a concrete type
	// s1 (Shape) -> (Rectangle)
	var s2 Rectangle = s1.(Rectangle)
	fmt.Println("Perimeter of s2 :", s2.Perimeter())
}
