/*
	Go does not support inheritance like traditional object-oriented languages, but has its own way - COMPOSITION.
	Inheritance is achieved by constructing anonymous properties in struct | embedding mechanism to create relationships between structs.
*/

package main

import "fmt"

type Contract struct {
	firstName   string
	lastName    string
	phoneNumber string
}

type Business struct {
	name    string
	address string
	// Anonymous property
	// Embed struct
	Contract
}

func (b Business) info() {
	// Access attributes directly or by anonymous property
	// fmt.Printf("Contract at %s is %s %s", b.name, b.Contract.firstName, b.Contract.lastName)
	fmt.Printf("Contract at %s is %s %s", b.name, b.firstName, b.lastName)
}

func main() {
	con := Contract{
		firstName:   "Kirito",
		lastName:    "Shiba",
		phoneNumber: "1234-1234",
	}

	bus := Business{
		name:     "Shop ABC",
		address:  "123 ZZ st",
		Contract: con,
	}
	bus.info()
}
