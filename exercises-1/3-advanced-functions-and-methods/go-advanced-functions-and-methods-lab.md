

# Go


## Advanced Functions and Closures

In this lab we will explore functions in Go and their use of memory. We will begin by examining several library functions, then we will look at their source code, finally applying what we see to our own code.


### 1. Go builtins

While most code we use or build in Go resides in a package.  The core language has several variables, functions, and types predefined, known as "builtins".  We will use the official documentation of these assorted builtins here https://golang.org/pkg/builtin/ or via `go doc builtin`.

Let's begin by taking a look at the `make()` function documentation.


#### func make(t Type, size ...IntegerType) Type

Per the documentation:

```
    The make built-in function allocates and initializes an object of type
    slice, map, or chan (only). Like new, the first argument is a type, not a
    value. Unlike new, make's return type is the same as the type of its
    argument, not a pointer to it. The specification of the result depends on
    the type:

        Slice: The size specifies the length. The capacity of the slice is
        equal to its length. A second integer argument may be provided to
        specify a different capacity; it must be no smaller than the
        length. For example, make([]int, 0, 10) allocates an underlying array
        of size 10 and returns a slice of length 0 and capacity 10 that is
        backed by this underlying array.
    
	    Map: An empty map is allocated with enough space to hold the
        specified number of elements. The size may be omitted, in which case
        a small starting size is allocated.
    
	    Channel: The channel's buffer is initialized with the specified
        buffer capacity. If zero, or the size is omitted, the channel is
        unbuffered.
```

* Create three separate functions which returning an object of type Slice, Map, and Channel respectively using the builtin `make` function.

To begin, create working directories for you experiments. Create a root directory for your executables and a `mylib` directory to house the mylib package:

```
~$ cd ~

~$ mkdir -p $(go env GOPATH)/src/lab-advanced-functions-and-methods/mylib

~$ cd $(go env GOPATH)/src/lab-advanced-functions-and-methods

~/go/src/lab-advanced-functions-and-methods$ ls -l

total 4
drwxr-xr-x 2 ubuntu ubuntu 4096 Jun  1 03:48 mylib

~/go/src/lab-advanced-functions-and-methods$
```

Create the following mylib.go library source file in the mylib directory:

```go
~/go/src/lab-advanced-functions-and-methods$ vi ./mylib/mylib.go

package mylib

func MakeSlice() []int {
	return make([]int, 10)
}

func MakeMap() map[string]int {
	return make(map[string]int, 10)
}

func MakeChannel() chan int {
	return make(chan int)
}

~/go/src/lab-advanced-functions-and-methods$
```

Create the following main.go source file in the cmd directory:

```go
~/go/src/lab-advanced-functions-and-methods$ vi main.go

package main

import (
	"fmt"
	"example.com/mylib"
)

func main() {
	s := mylib.MakeSlice()
	m := mylib.MakeMap()
	c := mylib.MakeChannel()

	fmt.Printf("Type: %T Value: %v\n", s, s)
	fmt.Printf("Type: %T Value: %v\n", m, m)
	fmt.Printf("Type: %T Value: %v\n", c, c)
}

~/go/src/lab-advanced-functions-and-methods$
```

Create the module.

```
~/go/src/lab-advanced-functions-and-methods$ go mod init example.com

go: creating new go.mod: module example.com
go: to add module requirements and sums:
	go mod tidy

~/go/src/lab-advanced-functions-and-methods$
```

And now, run the program.

```
~/go/src/lab-advanced-functions-and-methods$ go run main.go

Type: []int Value: [0 0 0 0 0 0 0 0 0 0]
Type: map[string]int Value: map[]
Type: chan int Value: 0xc000080060

~/go/src/lab-advanced-functions-and-methods$
```

* Is make a variadic function?


#### func new(Type) \*Type

Per the documentation:

```
    The new built-in function allocates memory. The first argument is a type,
    not a value, and the value returned is a pointer to a newly allocated zero
    value of that type.
```

* Create a series of functions using the `new` function.

Here is an example of using new with type int.

```go
~/go/src/lab-advanced-functions-and-methods$ vi ./mylib/new.go

package mylib

func NewInt() *int {
	return new(int)
}

func NewSlice() *[]int {
	return new([]int)
}

~/go/src/lab-advanced-functions-and-methods$
```

Create or modify the command to run our new wrapper function(s):

```go
~/go/src/lab-advanced-functions-and-methods$ vi main.go

package main

import (
	"fmt"
	"example.com/mylib"
)

func main() {
	s := mylib.MakeSlice()
	m := mylib.MakeMap()
	c := mylib.MakeChannel()

	fmt.Printf("--Make------------\n")
	fmt.Printf("Type: %T Value: %v\n", s, s)
	fmt.Printf("Type: %T Value: %v\n", m, m)
	fmt.Printf("Type: %T Value: %v\n", c, c)

	i := mylib.NewInt()
	ns := mylib.NewSlice()

	fmt.Printf("--New-------------\n")
	fmt.Printf("Type: %T Value: %v\n", i, i)
	fmt.Printf("Type: %T Value: %v\n", ns, ns)
}

~/go/src/lab-advanced-functions-and-methods$
```

Now run the program.

```
~/go/src/lab-advanced-functions-and-methods$ go run main.go

--Make------------
Type: []int Value: [0 0 0 0 0 0 0 0 0 0]
Type: map[string]int Value: map[]
Type: chan int Value: 0xc00005e060
--New-------------
Type: *int Value: 0xc00001a0b8
Type: *[]int Value: &[]

~/go/src/lab-advanced-functions-and-methods$
```


#### Comparing make and new

Now that we have a passing familiarity with `make` and `new`, what is going on, how do they compare? Both return an allocated place in memory, but only `make` initializes an object of the desired type.  The Go community has discussed combining the two functions into one. One benefit of `new` is its ability to return a pointer to a non-composite type (like an int).

For example, if you wanted to return a pointer to an int without using `new`, you would write the following:

```
func myNew() *int {
  var i int
  return &i
}
```

versus

```
i = new(int)
```


### 2. Reviewing the Source

With a couple of functions in play, lets look under the hood.

The current version of Go's gc compiler is written in Go. Prior to Go 1.4, gc was written in C. A real benefit of this is we can read the compiler and standard library code if we know Go!


### Review `make` Source Code

Lets first locate the code to the `make` function.  We know from earlier that the builtin types and functions are documented via a package called "builtin".  This gives us a clue on where to look in the source code.

If you already have a local copy of the source you can examine the code there, or you can clone the Go source from github, or browse the code on Github directly. To clone the source:

```
~/go/src/lab-advanced-functions-and-methods$ cd ~

# if git is not installed
~$ sudo apt install git -y

~$ git clone --depth 1 --branch go1.18.2 http://github.com/golang/go go1.18.2

...

~$
```

Now take a look at the code:

```
user@ubuntu:~$ cd ~/go1.18.2/src/ # or cd $(go env GOROOT/src)

~/go1.18.2/src$ grep -B 17 "func make" builtin/builtin.go

// The make built-in function allocates and initializes an object of type
// slice, map, or chan (only). Like new, the first argument is a type, not a
// value. Unlike new, make's return type is the same as the type of its
// argument, not a pointer to it. The specification of the result depends on
// the type:
//	Slice: The size specifies the length. The capacity of the slice is
//	equal to its length. A second integer argument may be provided to
//	specify a different capacity; it must be no smaller than the
//	length. For example, make([]int, 0, 10) allocates an underlying array
//	of size 10 and returns a slice of length 0 and capacity 10 that is
//	backed by this underlying array.
//	Map: An empty map is allocated with enough space to hold the
//	specified number of elements. The size may be omitted, in which case
//	a small starting size is allocated.
//	Channel: The channel's buffer is initialized with the specified
//	buffer capacity. If zero, or the size is omitted, the channel is
//	unbuffered.
func make(t Type, size ...IntegerType) Type

~/go1.18.2/src$
```

The builtin package is used for documentation. The implementation is elsewhere. Remembering that builtins are predefined, gives us a clue where to find them. There must be concrete types somewhere as we didn't supply them.

```
~/go1.18.2/src$ grep -nr "func make[s,m,c]" runtime

runtime/slice.go:38:func makeslicecopy(et *_type, tolen int, fromlen int, from unsafe.Pointer) unsafe.Pointer {
runtime/slice.go:88:func makeslice(et *_type, len, cap int) unsafe.Pointer {
runtime/slice.go:106:func makeslice64(et *_type, len64, cap64 int64) unsafe.Pointer {
runtime/chan.go:64:func makechan64(t *chantype, size int64) *hchan {
runtime/chan.go:72:func makechan(t *chantype, size int) *hchan {
runtime/map.go:283:func makemap64(t *maptype, hint int64, h *hmap) *hmap {
runtime/map.go:293:func makemap_small() *hmap {
runtime/map.go:304:func makemap(t *maptype, hint int, h *hmap) *hmap {

~/go1.17.2/src$
```

From the previous code listing we can see which files contain the implementations.

* Take a few minutes to review `runtime/chan.go`
  * What is hchan?
* Take a few minutes to review `runtime/slice.go`
  * What does growslice do?
* Take a few minutes to review `runtime/map.go`
  * Does map use `goto`?


### 3. Memory Analysis

Go has a unique way of using memory when programs are running. In general, the memory model is implementation dependent. Even with that, we have ways to understand at a high level what is going on.

```
~/go1.18.2/src$ cd ~/go/src/lab-advanced-functions-and-methods/

~/go/src/lab-advanced-functions-and-methods$
```

From a tooling perspective, `go run` calls `go build` which calls `go tool compile`.  

* Review `go tool compile -help` and related `go tool link -help`
* Review `go help build`
* Review `go help run`

Knowing that and digging into the help menus, we see we can pass debugging type flags down the build chain.  

Here we pass to the build step the `-m` flag via the `-gcflags` argument. A bit circular, you can find the list of gcflags in the `go tool compile -help` output.

```
~/go/src/lab-advanced-functions-and-methods$ go tool compile -help |& grep -P "\-m\t"

  -m	print optimization decisions

~/go/src/lab-advanced-functions-and-methods$  
```  

```
~/go/src/lab-advanced-functions-and-methods$ go run -gcflags -m main.go

# command-line-arguments
./main.go:9:22: inlining call to mylib.MakeSlice
./main.go:10:20: inlining call to mylib.MakeMap
./main.go:11:24: inlining call to mylib.MakeChannel
./main.go:13:12: inlining call to fmt.Printf
./main.go:14:12: inlining call to fmt.Printf
./main.go:15:12: inlining call to fmt.Printf
./main.go:16:12: inlining call to fmt.Printf
./main.go:18:19: inlining call to mylib.NewInt
./main.go:19:22: inlining call to mylib.NewSlice
./main.go:21:12: inlining call to fmt.Printf
./main.go:22:12: inlining call to fmt.Printf
./main.go:23:12: inlining call to fmt.Printf
./main.go:9:22: make([]int, int(10)) escapes to heap
./main.go:10:20: make(map[string]int, int(10)) escapes to heap
./main.go:14:12: ... argument does not escape
./main.go:14:13: s escapes to heap
./main.go:14:13: s escapes to heap
./main.go:15:12: ... argument does not escape
./main.go:16:12: ... argument does not escape
./main.go:18:19: new(int) escapes to heap
./main.go:19:22: new([]int) escapes to heap
./main.go:22:12: ... argument does not escape
./main.go:23:12: ... argument does not escape
--Make------------
Type: []int Value: [0 0 0 0 0 0 0 0 0 0]
Type: map[string]int Value: map[]
Type: chan int Value: 0xc000094060
--New-------------
Type: *int Value: 0xc0000c0050
Type: *[]int Value: &[]

~/go/src/lab-advanced-functions-and-methods$
```

For now, the thing to notice, does memory escape, if it does it goes on the heap, if it doesn't it stays on the stack.

* Create an example where memory goes on the stack versus the heap

Here is another example, but in this case we will show the assembly code to see what is going on.

```go
~/go/src/lab-advanced-functions-and-methods$ vi memory.go

package main

type structType struct{}

func func1() *structType {
	var chunk *structType = new(structType)
	return chunk
}

func func2() structType {
	var chunk structType
	return chunk
}

type bigStruct struct {
	lots [1e4]float64
}

func func3() bigStruct {
	var chunk bigStruct
	return chunk
}

func main() {
	func1()
	func2()
	func3()
}

~/go/src/lab-advanced-functions-and-methods$
```

At a high level, what we are about to do is translate the code into assembly and see which functions call `new`. `new` places objects on the heap. Many programmers are familiar with stacks and heaps, and generally believe simple data types (ex. int) live on the stack.  Unlike many languages, Go tries to make more efficient use of memory by keeping large objects on the stack in some cases.

Lets run the previous code and review the assembly.

```
~/go/src/lab-advanced-functions-and-methods$ go tool compile -help |& grep -P "\-S\t"

  -S	print assembly listing

~/go/src/lab-advanced-functions-and-methods$
```

```
~/go/src/lab-advanced-functions-and-methods$ go tool compile -S memory.go > memory.asm

~/go/src/lab-advanced-functions-and-methods$ head -45 memory.asm

"".func1 STEXT nosplit size=8 args=0x0 locals=0x0 funcid=0x0 align=0x0
	0x0000 00000 (memory.go:5)	TEXT	"".func1(SB), NOSPLIT|ABIInternal, $0-0
	0x0000 00000 (memory.go:5)	FUNCDATA	$0, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
	0x0000 00000 (memory.go:5)	FUNCDATA	$1, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
	0x0000 00000 (memory.go:7)	LEAQ	runtime.zerobase(SB), AX
	0x0007 00007 (memory.go:7)	RET
	0x0000 48 8d 05 00 00 00 00 c3                          H.......
	rel 3+4 t=14 runtime.zerobase+0
"".func2 STEXT nosplit size=1 args=0x0 locals=0x0 funcid=0x0 align=0x0
	0x0000 00000 (memory.go:10)	TEXT	"".func2(SB), NOSPLIT|ABIInternal, $0-0
	0x0000 00000 (memory.go:10)	FUNCDATA	$0, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
	0x0000 00000 (memory.go:10)	FUNCDATA	$1, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
	0x0000 00000 (memory.go:12)	RET
	0x0000 c3                                               .
"".func3 STEXT nosplit size=29 args=0x13880 locals=0x0 funcid=0x0 align=0x0
	0x0000 00000 (memory.go:19)	TEXT	"".func3(SB), NOSPLIT|ABIInternal, $0-80000
	0x0000 00000 (memory.go:19)	FUNCDATA	$0, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
	0x0000 00000 (memory.go:19)	FUNCDATA	$1, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
	0x0000 00000 (memory.go:19)	LEAQ	"".~r0+8(SP), DI
	0x0005 00005 (memory.go:19)	MOVL	$10000, CX
	0x000a 00010 (memory.go:19)	XORL	AX, AX
	0x000c 00012 (memory.go:19)	REP
	0x000d 00013 (memory.go:19)	STOSQ
	0x000f 00015 (memory.go:21)	LEAQ	"".~r0+8(SP), DI
	0x0014 00020 (memory.go:21)	MOVL	$10000, CX
	0x0019 00025 (memory.go:21)	REP
	0x001a 00026 (memory.go:21)	STOSQ
	0x001c 00028 (memory.go:21)	RET
	0x0000 48 8d 7c 24 08 b9 10 27 00 00 31 c0 f3 48 ab 48  H.|$...'..1..H.H
	0x0010 8d 7c 24 08 b9 10 27 00 00 f3 48 ab c3           .|$...'...H..
"".main STEXT nosplit size=1 args=0x0 locals=0x0 funcid=0x0 align=0x0
	0x0000 00000 (memory.go:24)	TEXT	"".main(SB), NOSPLIT|ABIInternal, $0-0
	0x0000 00000 (memory.go:24)	FUNCDATA	$0, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
	0x0000 00000 (memory.go:24)	FUNCDATA	$1, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
	0x0000 00000 (memory.go:28)	RET
	0x0000 c3                                               .
type..eq.[10000]float64 STEXT dupok nosplit size=45 args=0x10 locals=0x0 funcid=0x0 align=0x0
	0x0000 00000 (<autogenerated>:1)	TEXT	type..eq.[10000]float64(SB), DUPOK|NOSPLIT|ABIInternal, $0-16
	0x0000 00000 (<autogenerated>:1)	FUNCDATA	$0, gclocals·dc9b0298814590ca3ffc3a889546fc8b(SB)
	0x0000 00000 (<autogenerated>:1)	FUNCDATA	$1, gclocals·69c1753bd5f81501d95132d08af04464(SB)
	0x0000 00000 (<autogenerated>:1)	FUNCDATA	$5, type..eq.[10000]float64.arginfo1(SB)
	0x0000 00000 (<autogenerated>:1)	FUNCDATA	$6, type..eq.[10000]float64.argliveinfo(SB)
	0x0000 00000 (<autogenerated>:1)	PCDATA	$3, $1
	0x0000 00000 (<autogenerated>:1)	XORL	CX, CX
	0x0002 00002 (<autogenerated>:1)	JMP	7

~/go/src/lab-advanced-functions-and-methods$
```

If you search the output for `newobject`, you will see calls that allocate heap. We don't see any.....

In the previous code we see `func1` allocating memory on the heap (we thought!). One interesting thing is look at `func3` Go source. Our bigStruct datatype has `1e4` elements of type `float64`.  That is pretty big, however there is no call to `newobject`, meaning its on the stack. If you are familiar with closures, you may wonder why its not on the heap as well, but Go has analyzed our code and knows there are no references to it, so it can live on the dynamically sized stack.

Lets force the issue, lets double the size of the array to `1e8` and rerun our example.

```
type bigStruct struct {
    lots [1e4]float64
}
```

becomes

```
type bigStruct struct {
    lots [1e8]float64
}
```

```
~/go/src/lab-advanced-functions-and-methods$ go tool compile -S memory.go > memory.asmv2

~/go/src/lab-advanced-functions-and-methods$
```

Review the new assembly listing.

Now if we search for calls (new), we see `func3` allocating to the heap.  Specifically in this labs example, line 20, which calls `var chunk bigStruct`, in turn in assembly is `runtime.newobject`.

```
~/go/src/lab-advanced-functions-and-methods$ grep new memory.asmv2 

	0x002a 00042 (memory.go:20)	CALL	runtime.newobject(SB)
	rel 43+4 t=7 runtime.newobject+0
	0x0020 00032 (memory.go:20)	CALL	runtime.newobject(SB)
	0x0031 00049 (memory.go:27)	CALL	runtime.newobject(SB)
	rel 33+4 t=7 runtime.newobject+0
	rel 50+4 t=7 runtime.newobject+0

~/go/src/lab-advanced-functions-and-methods$
```

I guess there is a heap after all!

* Run the memory optimizer gcflag "-m" on your memory.go source, does it look as expected?
* How about two -m? "go tool compile -S -m -m memory.go" (hint gives even more details)

> nb. Use the `diff -y f1.txt f2.txt` to easily see the differences between files


### 4. Additional Functions

Using our previous examples as a template:

* review the following builtin functions, then
* create a sample program with them, then
* review implementation code, then
* and check memory allocation.

for the following builtins:

* `cap`
* `len`
* `append`

`cap`, `len`, and `append` are used to check and resize a slice.


### 5. Challenge: Closures

In a file called "fib.go," implement a stateful Fibonacci function that uses a closure to return the next value in the Fibonacci sequence each time it is called.

> The Fibonacci Sequence is the series of numbers: 0, 1, 1, 2, 3, 5, 8, ... where the next number is found by adding the two prior numbers.

Use this as starter code:

```go
package main

import "fmt"

// fibonacci is a function that returns
// a function that returns an int.
func fibonacci() func() int {
}

func main() {
	f := fibonacci()
	for i := 0; i < 10; i++ {
		fmt.Println(f())
	}
}
```

- Can you generate two independent sequences using the fibonacci() function?
- Does the caller need to know anything about the intermediate state?

<br>

Congratulations you have completed the lab!

<br>

