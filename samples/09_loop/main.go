package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var pl = fmt.Println

func main() {
	// === FOR loop ===
	// for initialization; conditions; postStatement {BODY}
	for index := 1; index <= 5; index++ {
		pl(index)
	}

	for index := 5; index > 0; index-- {
		pl(index)
	}

	// Error: index out of scope
	// pl("index:", index)

	// === WHILE loop ===
	x := 0
	for x < 5 {
		pl(x)
		x++
	}

	// === INFINITE loop ===
	// for {}
	for true {
		randNum := randomInt(50)
		fmt.Printf("Guess a number between 0 and 50 (number is %d): ", randNum)

		var guess string
		fmt.Scan(&guess)
		guess = strings.TrimSpace(guess)
		guessNum, err := strconv.Atoi(guess)
		if err != nil {
			log.Fatal(err)
		}

		if guessNum > randNum {
			pl("Pick a lower value")
		} else if guessNum < randNum {
			pl("Pick a higher value")
		} else {
			pl("You guessed it")
			break
		}
	}
}

func randomInt(max int) int {
	seedSecs := time.Now().Unix()
	rand.Seed(seedSecs)
	return rand.Intn(max) + 1
}
