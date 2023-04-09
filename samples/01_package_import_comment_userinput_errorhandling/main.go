// package: a collection of code
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

// One line comment

/*
Multi-line comment
*/

func main() {

	// User input
	firstName := getFirstName()
	lastName, err := getLastName()

	// Blank identifier - ignore unused variables
	// Ignore error
	// name, _ := reader.ReadString('\n')

	// Error handling
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Hello", firstName, lastName)
	}
}

// getFirstName: get user input using fmt.Scan
func getFirstName() string {
	var firstName string

	fmt.Print("Enter your first name: ")
	fmt.Scan(&firstName)
	return firstName
}

// getLastName: get user input using bufio.ReadString
func getLastName() (string, error) {
	fmt.Print("Enter your last name: ")
	reader := bufio.NewReader(os.Stdin)
	lastName, err := reader.ReadString('\n')
	return lastName, err
}
