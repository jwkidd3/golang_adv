package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"runtime/trace"
)

func nogo(c chan int, id int) {
	i := 0
	alot := 100000
	randomMath := 0
	for i < alot {
		randomMath = i + randomMath
		i++
	}

	c <- id
}

func slowgo(c chan int, id int) {
	s := make([]int, rand.Intn(10000))

	randomMath := 0

	for i, v := range s {
		randomMath = i*v + randomMath
	}

	c <- id
}

func fastgo(c chan int, id int) {
	c <- id
}

func cpu(count int) {
	c := make(chan int)

	id := 0

	incrId := func() int {
		x := id
		id++
		return x
	}

	decrId := func() int {
		x := id
		id--
		return x
	}

	for id < count {
		go nogo(c, id)
		go slowgo(c, incrId())
		go fastgo(c, incrId())
	}

	for id >= 0 {
		fmt.Print("\033[H\033[2J")
		fmt.Printf("\033[0;0H")
		fmt.Println("got", <-c, id)
		decrId()
	}

	fmt.Println("finished")
}

var count = flag.Int("count", 1000, "count")

func main() {
	f, err := os.Create("trace.out")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = trace.Start(f)
	if err != nil {
		panic(err)
	}
	defer trace.Stop()
	flag.Parse()

	cpu(*count)
}

