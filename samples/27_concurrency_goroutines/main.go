/*
	CONCURRENCY

	The ability of a program to execute multiple tasks or processes simultaneously, independently, and in parallel,
	Concurrency is achieved through the use of goroutines and channels.

	A GOROUTINE is a lightweight thread of execution managed by the Go runtime that runs concurrently with other goroutines within the same address space.
*/

package main

import (
	"fmt"
	"time"
)

func printTo5() {
	for index := 0; index < 5; index++ {
		fmt.Println("Func 1 :", index)
	}
}

func printTo10() {
	for index := 0; index < 10; index++ {
		fmt.Println("Func 2 :", index)
	}
}

func main() {
	go printTo5()
	go printTo10()

	time.Sleep(2 * time.Millisecond)
}
