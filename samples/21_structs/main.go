/* 
	STRUCT
	
	Allow to store multiple values with different data types in structured way
*/

package main

import "fmt"

type Customer struct {
	fullName string
	address  string
	balance  float64
}

// Function work with a struct
func getInfo(c Customer) {
	fmt.Printf("%s (%s) owes us %.2f\n", c.fullName, c.address, c.balance)
}

func addNewAddress(c *Customer, address string) {
	c.address = address
}

func main() {
	// Declare a struct
	var ks Customer
	ks.fullName = "Kirito Shiba"
	ks.address = "10 ABC st"
	ks.balance = 3200.99

	getInfo(ks)
	addNewAddress(&ks, "12 DEF st")
	fmt.Printf("Now %s live at %s\n", ks.fullName, ks.address)

	// Define struct with field names
	ay := Customer{
		fullName: "Asuna Yuuki",
		address:  "12 DEF st",
		balance:  40000,
	}
	getInfo(ay)

	// Define struct with ordered fields
	ksu := Customer{"Kirigaya Suguha", "10 ABC st", 256.789}
	getInfo(ksu)
}
