package main

import "fmt"

/*
	Maps

A map maps keys to values.

The zero value of a map is nil. A nil map has no keys, nor can keys be added.

The make function returns a map of the given type, initialized and ready for use.

	var myMap map [keyType]valueType
	myMap := make(map[keyType]valueType)
*/
var pl = fmt.Println

func main() {
	// Declare a map variable
	var heroes map[string]string
	// Create the map
	heroes = make(map[string]string)

	// Add keys, values
	heroes["Iron Man"] = "Tony Stark"
	heroes["Spider-Man"] = "Peter Parker"
	heroes["Captain America"] = "Steve Rogers"
	heroes["Hulk"] = "Bruce Banner"

	// Define map in 1 step
	villains := make(map[string]string)
	villains["Loki"] = "Tom Hiddleston"
	villains["Venom"] = "Eddie Brock"

	// Define with map literal
	heroWifes := map[int]string{
		1: "Pepper Potts",
		2: "Peggy Carter",
	}

	// Get the map values
	// using %v because of not always know the value data type
	fmt.Printf("Iron Man is %v\n", heroes["Iron Man"])
	pl("Iron Man's wife :", heroWifes[1])

	// Access a not existed key => get nil
	pl("Some wife :", heroWifes[3])

	// Check if there is a value or nil
	_, ok := heroWifes[3]
	pl("Is there a 3rd wife :", ok)

	// Iterate over the map
	for key, value := range heroes {
		fmt.Printf("%s is %s\n", key, value)
	}

	// Delete a key value pair
	delete(heroes, "Hulk")

	pl("Heroes :", heroes)
}
