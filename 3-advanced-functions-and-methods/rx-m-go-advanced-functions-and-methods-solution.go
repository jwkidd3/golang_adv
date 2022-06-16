```
package main

import "fmt"

// fibonacci is a function that returns
// a function that returns an int.
func fibonacci() func() int {
	f1 := 1
	f2 := 0
	return func() int {
		f2, f1 = (f1 + f2), f2
		return f1
	}
}

func main() {
	f := fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Println(f())
	}
}
```

```
go run fib.go 
0
1
1
2
3
5
8
13
21
34
```