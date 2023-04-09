

# Go


## CSP and goroutines

Per the CSP documentation (https://go.dev/doc/effective_go#concurrency):

> Concurrent programming is a large topic and there is space only for some Go-specific highlights here.
> 
> Concurrent programming in many environments is made difficult by the subtleties required to implement correct access to shared variables. Go encourages a different approach in which shared values are passed around on channels and, in fact, never actively shared by separate threads of execution. Only one goroutine has access to the value at any given time. Data races cannot occur, by design. To encourage this way of thinking we have reduced it to a slogan:
> 

> &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Do not communicate by sharing memory; instead, share memory by communicating.

>
> This approach can be taken too far. Reference counts may be best done by putting a mutex around an integer variable, for instance. But as a high-level approach, using channels to control access makes it easier to write clear, correct programs.
> 
> One way to think about this model is to consider a typical single-threaded program running on one CPU. It has no need for synchronization primitives. Now run another such instance; it too needs no synchronization. Now let those two communicate; if the communication is the synchronizer, there's still no need for other synchronization. Unix pipelines, for example, fit this model perfectly. Although Go's approach to concurrency originates in Hoare's Communicating Sequential Processes (CSP), it can also be seen as a type-safe generalization of Unix pipes.


### 1. goroutines

To run code as a goroutine, we prefix the call with `go`.  Here is an example where a long running function normally taking 10 seconds, so on the second execution we run it as a goroutine.

```
~$ mkdir -p $(go env GOPATH)/src/lab-csp-and-go-routines && cd $_

~/go/src/lab-csp-and-go-routines$ vi lab-csp-a.go

package main

import (
	"fmt"
	"time"
)

func main() {
	f := func() {
		fmt.Println("before")
		time.Sleep(10 * time.Second)
		fmt.Println("after")
	}

	f()

	go f()
}

~/go/src/lab-csp-and-go-routines$
```

Run our sample program.

```
~/go/src/lab-csp-and-go-routines$ go run lab-csp-a.go

before
after

~/go/src/lab-csp-and-go-routines$
```

It seems things executed but we don't have any output associated with the goroutine invocation of the function `f`. Another issue is there is no way for the goroutine to signal its completion to the caller. Unlike threads, there is no associated id (tid).


### 2. Channels

In order to communicate between the coroutine and the caller we leverage a channel.

To create a channel we use `make`.  Here we modify our existing program to leverage a channel to return the result from the function being executed in the goroutine.

We modify the second call to `f()` by wrapping it in another anonymous function and executing it as a goroutine. A closure is created with the channel `c`. We also add a print statement to show activity beyond the goroutine and ultimately we wait at the end to read from the channel `c` before printing a final time the return value.

```go
~/go/src/lab-csp-and-go-routines$ vi lab-csp-b.go

package main

import (
	"fmt"
	"time"
)

func main() {
	f := func() int {
		fmt.Println("f:before")
		time.Sleep(10 * time.Second)
		fmt.Println("f:after")
		return 0
	}

	f()

	c := make(chan int)

	go func() {
		fmt.Println("triggering f via go routine")
		c <- f()
	}()

	fmt.Println("main:do some other work before we block on <-")

	returnVal := <-c

	fmt.Println("main:return value of", returnVal)
}

~/go/src/lab-csp-and-go-routines$
```

Run the program.

```
~/go/src/lab-csp-and-go-routines$ go run lab-csp-b.go

f:before
f:after
main:do some other work before we block on <-
triggering f via go routine
f:before
f:after
main:return value of 0

~/go/src/lab-csp-and-go-routines$
```

We see in the previous example some basic usage of the channel, the key is to notice how things are being synchronized without the use of locks. It may not be clear if you have not used locks in your own code before.

Lets look at another example where we try simulate using all the CPUs available to us. Here is one way to see how many CPUs you have (under Linux):

```
~/go/src/lab-csp-and-go-routines$ lscpu

Architecture:                    x86_64
CPU op-mode(s):                  32-bit, 64-bit
Byte Order:                      Little Endian
Address sizes:                   46 bits physical, 48 bits virtual
CPU(s):                          2
On-line CPU(s) list:             0,1
Thread(s) per core:              2
Core(s) per socket:              1
Socket(s):                       1
NUMA node(s):                    1
Vendor ID:                       GenuineIntel
CPU family:                      6
Model:                           85
Model name:                      Intel(R) Xeon(R) Platinum 8175M CPU @ 2.50GHz
Stepping:                        4
CPU MHz:                         2499.998
BogoMIPS:                        4999.99
Hypervisor vendor:               KVM
Virtualization type:             full
L1d cache:                       32 KiB
L1i cache:                       32 KiB
L2 cache:                        1 MiB
L3 cache:                        33 MiB
NUMA node0 CPU(s):               0,1
Vulnerability Itlb multihit:     KVM: Mitigation: VMX unsupported
Vulnerability L1tf:              Mitigation; PTE Inversion
Vulnerability Mds:               Vulnerable: Clear CPU buffers attempted, no microcode; SMT Host state unk
                                 nown
Vulnerability Meltdown:          Mitigation; PTI
Vulnerability Spec store bypass: Vulnerable
Vulnerability Spectre v1:        Mitigation; usercopy/swapgs barriers and __user pointer sanitization
Vulnerability Spectre v2:        Mitigation; Full generic retpoline, STIBP disabled, RSB filling
Vulnerability Srbds:             Not affected
Vulnerability Tsx async abort:   Vulnerable: Clear CPU buffers attempted, no microcode; SMT Host state unk
                                 nown
Flags:                           fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 c
                                 lflush mmx fxsr sse sse2 ss ht syscall nx pdpe1gb rdtscp lm constant_tsc 
                                 rep_good nopl xtopology nonstop_tsc cpuid tsc_known_freq pni pclmulqdq ss
                                 se3 fma cx16 pcid sse4_1 sse4_2 x2apic movbe popcnt tsc_deadline_timer ae
                                 s xsave avx f16c rdrand hypervisor lahf_lm abm 3dnowprefetch invpcid_sing
                                 le pti fsgsbase tsc_adjust bmi1 hle avx2 smep bmi2 erms invpcid rtm mpx a
                                 vx512f avx512dq rdseed adx smap clflushopt clwb avx512cd avx512bw avx512v
                                 l xsaveopt xsavec xgetbv1 xsaves ida arat pku ospke

~/go/src/lab-csp-and-go-routines$
```

The key field of course is the **CPU(s):** count.

```
~/go/src/lab-csp-and-go-routines$ lscpu | grep ^CPU\(s\): | awk '{print $2}'

2

~/go/src/lab-csp-and-go-routines$
```

Now lets create a program that uses the number of CPUs to run some code, and then wait for all of them to complete in the driver (main) program.

```go
~/go/src/lab-csp-and-go-routines$ vi lab-csp-c.go

package main

import (
	"fmt"
	"runtime"
	"time"
)

var numCPU = runtime.NumCPU()

// Apply the operation to v[i], v[i+1] ... up to v[n-1].
func DoSome(cpu int, c chan int) {
	fmt.Println("Doing something on", cpu)
	time.Sleep(2 * time.Second)
	c <- 1 // signal that this piece is done
}

func DoAll() {
	c := make(chan int, numCPU) // Buffering optional but sensible.

	for i := 1; i <= numCPU; i++ {
		go DoSome(i, c)
	}
	// Drain the channel.
	for i := 0; i < numCPU; i++ {
		<-c // wait for one task to complete
	}
	// All done.
	fmt.Println("I only run once all the farmed out CPU work is done")
}

func main() {
	DoAll()
}

~/go/src/lab-csp-and-go-routines$
```

In the function `DoAll()`, we create a channel (c), then we create a goroutine that executes `DoSome()` and pass the channel to it. These means we can pass channels as first class citizens.  In each instance of the goroutine, we do some work and send a result back to the channel. Back in the `DoAll()` function, we wait for the same number of results to return via the channel.  In our case, we don't what the result is, we just want to make sure we get the same count. Finally, there is a Println that is represents further computation that was blocked until all the *numCPU* goroutines finished.

With that understanding, take your program for a run.

```
~/go/src/lab-csp-and-go-routines$ go run lab-csp-c.go

Doing something on 2
Doing something on 1
I only run once all the farmed out CPU work is done

~/go/src/lab-csp-and-go-routines$
```

* Should 2 print before 1?

There is more to goroutines and channels, but this hopefully gives a feel to doing synchronization without locking.


### 3. Challenge

Create a simple ping pong game using goroutines and channels. Two players (goroutines) send the ball back and forth via a channel (table). Create a function play() that takes two ints - the chance out of 100 that each player will hit the ball successfully on a given round (e.g. a highly skilled player might have 90). These probabilities should be used to determine if a given player hits the ball in the function hits(). Set up a loop that runs until a player fails to hit the ball. Print out a count each iteration and a message declaring who won and who lost. Use fmt.Scanln() so that the user must press enter between turns and to start the game. 

Here is some starter code:

```go
package main

import "fmt"
import "math/rand"
import "time"

func hits(probability int) bool {
	// given a probability, use randomness to determine if the player hits the ball
	// go here: https://gobyexample.com/random-numbers to learn more about the random package
}

func play(player0prob, player1prob int) {
	// make a channel
	fmt.Println("Press enter to begin")
	var input string
	fmt.Scanln(&input)
	playing := true
	// initialize other useful variables
    	for playing {
		// print out a rally count
		// check if the player hit the ball
		// if they did, send a message through the channel and print it out
		// make the user press enter before continuing to player1's turn
    	}
	// print out who won

}
func main() {
    	play(94, 66)
    	// a player that is very good and a player that is okay
}
```

<br>

Congratulations you have completed the lab!

<br>

