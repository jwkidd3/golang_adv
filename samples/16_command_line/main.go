package main

import (
	"fmt"
	"os"
	"strconv"
)

func main() {
	fmt.Println(os.Args)

	// Get all values after the first index (exclude executing command)
	args := os.Args[1:] // []string

	// Create []int from []string
	iArgs := []int{}
	for _, v := range args {
		val, err := strconv.Atoi(v)
		if err != nil {
			panic(err)
		}
		iArgs = append(iArgs, val)
	}

	// Find the max value
	max := 0
	for _, val := range iArgs {
		if val > max {
			max = val
		}
	}

	fmt.Println("Max value :", max)
}
