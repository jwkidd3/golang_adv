package main

import (
	"fmt"
	"reflect"
	"strconv"
)

var pl = fmt.Println

func main() {
	/* Data types */
	// 				int	|	float64	|	bool	|	string	|	rune
	// Default value: 0 |	0.0		|	false	|	""		|	0
	a := 12
	b := 3.14
	c := true
	d := "Hello"
	e := '😆'
	var f rune = '😆'

	/* Get variable data type */
	pl("a -", reflect.TypeOf(a))
	pl("b -", reflect.TypeOf(b))
	pl("c -", reflect.TypeOf(c))
	pl("d -", reflect.TypeOf(d))
	pl("e -", reflect.TypeOf(e))
	pl("f -", reflect.TypeOf(f))

	/* Cast data type */
	// Convert float to int
	cV1 := 2.52
	cV2 := int(cV1)
	pl("cV1 -", cV1, "- cV2 -", cV2)

	// Convert string to int
	cV3 := "2000"
	cV4, err := strconv.Atoi(cV3)
	pl("cV4 -", cV4, err, reflect.TypeOf(cV4))

	// Convert int to string
	cV5 := 2000
	cV6 := strconv.Itoa(cV5)
	pl("cV6 -", cV6, reflect.TypeOf(cV6))

	// Convert float to string
	cV7 := 2.52
	cV8 := strconv.FormatFloat(cV7, 'f', -1, 64)
	pl("cV8 -", cV8, reflect.TypeOf(cV8))

	cV9 := fmt.Sprintf("%f", cV7)
	pl("cV9 -", cV9, reflect.TypeOf(cV9))

	// Convert string to float
	cV10 := "3.14"
	if cV11, err := strconv.ParseFloat(cV10, 64); err == nil {
		pl("cV11 -", cV11, reflect.TypeOf(cV11))
	}
}
