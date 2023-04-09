/*
	CHANNELS

	are the pipes that connect concurrent goroutines.
	You can send values into channels from one goroutine and receive those values into another goroutine.

	A buffered channel is a type of channel that can hold a specified number of values in a buffer (channel with fixed size of the buffer)
*/

package main

import "fmt"

var pl = fmt.Println

func nums1(channel chan int) {
	channel <- 1
	channel <- 2
	channel <- 3
	// channel <- 4 // fatal error: all goroutines are asleep - deadlock! (goroutine 1 [chan receive (len: 3)])
}

func nums2(channel chan int) {
	channel <- 4
	channel <- 5
	channel <- 6
}

func getOrderedResults() {
	chan1 := make(chan int, 3)
	chan2 := make(chan int, 3)

	go nums1(chan1)
	go nums2(chan2)

	pl(<-chan1)
	pl(<-chan1)
	pl(<-chan1)

	pl(<-chan2)
	pl(<-chan2)
	pl(<-chan2)
}

func fibonancci(n int, c chan int) {
	x, y := 0, 1
	for index := 0; index < n; index++ {
		c <- x
		x, y = y, x+y
	}

	// Close channel
	close(c)
}

func main() {
	// Always get the same ordered results
	getOrderedResults()

	// Send these values into the channel with a corresponding concurrent receive.
	chan1 := make(chan string)
	go func() { chan1 <- "Hello" }()
	go func() { chan1 <- "World" }()

	a, b := <-chan1, <-chan1
	pl(a, b)

	// Buffered channel
	chan2 := make(chan int, 10)
	fibonancci(cap(chan2), chan2)

	// range c will receive values from channel chan2 until it is closed
	for index := range chan2 {
		pl(index)
	}
}
