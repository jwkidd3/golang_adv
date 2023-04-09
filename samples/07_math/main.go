package main

import (
	"fmt"
	"math"
)

var pl = fmt.Println

func main() {
	// Common operators
	pl("4 + 3 =", 4+3)
	pl("4 - 3 =", 4-3)
	pl("4 * 3 =", 4*3)
	pl("4 / 3 =", 4/3)
	pl("4 % 3 =", 4%3)

	counter := 1
	counter = counter + 1
	counter += 1
	counter++

	// Math functions
	pl("Min(4, 3) =", math.Min(4, 3))
	pl("Max(4, 3) =", math.Max(4, 3))
	pl("Abs(-20) =", math.Abs(-20))
	pl("Pow(2, 3) =", math.Pow(2, 3))
	pl("Sqrt(64) =", math.Sqrt(64))
	pl("Ceil(4.4) =", math.Ceil(4.4))
	pl("Floor(4.4) =", math.Floor(4.4))
	pl("Round(4.4) =", math.Round(4.4))
	pl("Log2(8) =", math.Log2(8))
	pl("Log10(100) =", math.Log10(100))
	pl("Exp(2) =", math.Exp(2))
	pl("Log(7.389) =", math.Log(math.Exp2(2))) // Log of e to the power of 2

	// Convert radians to degrees
	r90 := 90 * math.Pi / 180
	d90 := r90 * (180 / math.Pi)
	fmt.Printf("%.2f radians = %.2f degrees\n", r90, d90)
	pl("Sin(90) =", math.Sin(r90))
}
