/*
	Getters/Setters

	Encapsulation (in OOP):
		Getter and setter methods help to encapsulate the data within an object
		making it more secure and easier to maintain, control over how data is accessed and modified
*/

package main

import (
	"fmt"
	"log"

	model "example.com/project/models"
)

func main() {
	date := model.Date{}
	err := date.SetDay(20)
	if err != nil {
		log.Fatal(err)
	}

	err = date.SetMonth(10)
	if err != nil {
		log.Fatal(err)
	}

	err = date.SetYear(1990)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Date : %d/%d/%d\n", date.Year(), date.Month(), date.Day())
}
