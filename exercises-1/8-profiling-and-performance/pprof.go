package main

import "fmt"
import "runtime/pprof"
import "log"
import "os"

func bigBytes() *[]byte {
	s := make([]byte, 1000000000)
	return &s
}

func mem(iter int) int {
	total := 0
	i := 0

	for i <= iter {
		total += 1
		i++
		s := bigBytes()
		fmt.Println("mem", i, iter)
		if s == nil {
			fmt.Println("trouble")
		}
	}

	f, err := os.Create("./mem.proc")
	if err != nil {
		log.Fatal(err)
	}
	pprof.WriteHeapProfile(f)
	f.Close()

	return total
}

func main() {
	mem(10)
}
