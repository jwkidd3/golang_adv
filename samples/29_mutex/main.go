/*
	sync.Mutex

	A MUTEX (short for "mutual exclusion") is a synchronization primitive used to protect shared resources (such as variables, data structures, or I/O devices) from concurrent access by multiple goroutines at a time to avoid conflicts
	Support 2 methods:
		Lock() - Acquires the mutex and blocks until it is available. If the mutex is already locked by another goroutine, Lock() will block until it becomes available.
		Unlock() - Releases the mutex, allowing other goroutines to acquire it.
*/

package main

import (
	"fmt"
	"sync"
	"time"
)

type Account struct {
	balance int
	mu      sync.Mutex
}

func (a *Account) GetBalance() int {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.balance
}

func (a *Account) Withdraw(val int) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if val > a.balance {
		fmt.Println("Your withdrawal request has been declined due to insufficient funds")
	} else {
		a.balance -= val
		fmt.Printf("%d withdrawn - Balance : %d\n", val, a.balance)
	}
}

func main() {
	var acc Account
	acc.balance = 1000
	fmt.Println("Balance :", acc.GetBalance())

	for index := 0; index < 20; index++ {
		go acc.Withdraw(150)
	}

	// Wait for all the goroutines to complete by sleeping for one second
	time.Sleep(time.Second)
}
