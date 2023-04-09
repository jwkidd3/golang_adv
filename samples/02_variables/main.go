/*
	VARIABLES
		var name type	| var name type = value	| name := value
		- name: begin with letter, may contains letters, digits
		- Variable is a mutable data type, you can change its value but cannot change its type

	GLOBAL LEVEL VARIABLES
		var GlobalVar = "somevalue"
		Variable name with Uppercase first letter - exported
		Can be used everywhere across all packages

	PACKAGE LEVEL VARIABLES
		var aVar = "somevalue"
		Can be used everywhere in the same package
		All the functions in package can access to these variables
*/

package main

import "fmt"

var pl = fmt.Println

func main() {
	var vName string = "Hoang"
	var v1, v2 int = 1, 2
	var v3 = 3.14

	v4 := "Hello"
	v4 = "World"

	pl(vName, v1, v2, v3, v4)
}
