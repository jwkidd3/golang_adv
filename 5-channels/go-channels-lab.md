
# Go


## Channels as first class values

In this lab, we will explore channels at a deeper level. Recall that channels provide a way for two goroutines to
communicate with one another and synchronize their execution.

```
~$ mkdir -p $(go env GOPATH)/src/lab-channels-as-first-class-values && cd $_

~/go/src/lab-channels-as-first-class-values$
```


### 1. Unbuffered

As the Go [documentation](https://golang.org/doc/effective_go.html#channels) states: Unbuffered channels combine
communication — the exchange of a value — with synchronization — guaranteeing that two calculations (goroutines) are in
a known state. Receivers always block until there is data to receive. If the channel is unbuffered, the sender blocks
until the receiver has received the value. If the channel has a buffer, the sender blocks only until the value has been
copied to the buffer; if the buffer is full, this means waiting until some receiver has retrieved a value. Let's take a
look at a basic example.

```go
~/go/src/lab-channels-as-first-class-values$ vi main.go

package main

import "fmt"

func main() {

  var c = make(chan int)

  go func() { c<-1 }()

  v := <-c

  fmt.Println("done:", v)
}

~/go/src/lab-channels-as-first-class-values$
```

Run the program.

```
~/go/src/lab-channels-as-first-class-values$ go run main.go

done: 1

~/go/src/lab-channels-as-first-class-values$
```

* What happens when you comment out the goroutine line?


### 2. Pipelining

The Go [documentation](https://go.dev/blog/pipelines) states that there's no formal definition of a pipeline in Go;
it's just one of many kinds of concurrent programs. Informally, a pipeline is a series of stages connected by channels,
where each stage is a group of goroutines running the same function.

In each stage, the goroutines:

	* receive values from upstream via inbound channel(s)
	* perform some function on that data, usually producing new values
	* send values downstream via outbound channel(s)

Each stage has any number of inbound and outbound channels, except the first and last stages, which have only outbound
or inbound channels, respectively. The first stage is sometimes called the source or producer; the last stage, the sink
or consumer.

```go
~/go/src/lab-channels-as-first-class-values$ vi main.go

package main

import "fmt"

func main() {

	var c = make(chan int)
	var c2 = make(chan string)

	go func() { c <- 1 }()
	go func() {
		<-c
		c2 <- "ok"
	}()

	v := <-c2

	fmt.Println("done:", v)
}

~/go/src/lab-channels-as-first-class-values$
```

```
~/go/src/lab-channels-as-first-class-values$ go run main.go

done: ok

~/go/src/lab-channels-as-first-class-values$
```

* Add a call to time.Sleep between reading the channel and writing the next channel
* How large a message can be placed on a channel?


### 3. Buffered

To specify a channel as buffered, provide the buffer length to make when initializing a channel. In the example below,
`c` only sends to a buffered channel block when the buffer is full (when it reaches 10 ints), and it receives block when
the buffer is empty.

```go
~/go/src/lab-channels-as-first-class-values$ vi main.go

package main

import (
	"fmt"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(1)
}

func randSleep() {
	s := time.Duration(rand.Intn(1000)) * time.Millisecond
	time.Sleep(s)
}

func display() {
	fmt.Printf("\rsend:%d\treceive:%d", s, r)
}

var s int
var r int

func main() {
	var c = make(chan int, 10)

	count := 30

	go func() {
		for i := 0; i < count; i++ {
			randSleep()
			c <- i
			s++
			display()
		}
	}()

	for k := 0; k < count; k++ {
		randSleep()
		<-c
		r++
		display()
	}

	display()

	fmt.Println("\nfinished")
}

~/go/src/lab-channels-as-first-class-values$
```

Run the program.

```
~/go/src/lab-channels-as-first-class-values$ go run main.go

30      30
finished

~/go/src/lab-channels-as-first-class-values$
```

* Do you think buffered or un-buffered channels would perform faster?


### 4. Select

What does this code do https://go.dev/play/p/Vco7d8Lmhn

* Modify the code to use `range` and `close`, versus a for loop on receiving.
* A potential design to use a single "for" loop (not common mind you) - https://rillabs.com/posts/range-over-multiple-go-channels


### 5. Unidirectional

We can at compile time prevent channels from being used in both directions.  We can do this by passing in the channel
versus accessing an externally accessed channel like the previous examples.

```go
...
func sendOnly(sendOn chan<-int, msg) { sendOn<-msg }

func receieveOnly(receive <-chan int) { value := <- receive }
...
```

* Create a program that only allows us to send on a channel and receive on another channel
* Try to read from a write only channel, what happens?

By passing channels as arguments we see they are first class values.

<br>

Congratulations you have completed the lab!

<br>

