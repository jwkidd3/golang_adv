
# Go


## Generics

Generics is a programming mechanism where data types are abstracted from algorithms. In a programming language where
generics are not available, you often have to duplicate code when using the same code with different paramter types. 
This lack of generics can lead to binary bloat, duplicated bugs, missed opitimziations and more.

Some languages, including Go, do not allow function overloading. This very issue is addressed in the Go FAQ:

https://go.dev/doc/faq#overloading

> Why does Go not support overloading of methods and operators?

> Method dispatch is simplified if it doesn't need to do type matching as well. Experience with other languages told us that having a variety of methods with the same name but different signatures was occasionally useful but that it could also be confusing and fragile in practice. Matching only by name and requiring consistency in the types was a major simplifying decision in Go's type system.

> Regarding operator overloading, it seems more a convenience than an absolute requirement. Again, things are simpler without it.

Go produces the following error if overloading is attemped.

```
# sample Go code
func A(a int64) int64 { return a }
func A(a float64) float64 { return a}

# go run code.go
# command-line-arguments
./overloaded.go:5:6: A redeclared in this block
	./overloaded.go:6:6: other declaration of A
```

A simple solution is to change the second function name to "B". A less simple solution would be to abstract the 
parameter (eg. using `any` type as an example) and have one function (eg. casting as needed from `any`, or using 
reflection). The former would also impose documentation concerns and the later seems to broad in power.


### 1. Filtering A Collection Without Generics

We now show an example where two filter functions are used to process arrays of a particular type. Prior to Go 1.18
we are required to have two functions unless we leverage the `any` (the empty interface).

```
~$ cd

~$ mkdir ./generics

~$ cd ./generics

~/generics$ vi no_generics.go

package main

import (
	"fmt"
)

func FilterInt(s []int, f func(int) bool) []int {
	var r[]int
	for _, v := range s {
		if f(v) {
			r = append(r, v)
		}
	}

	return r
}

func FilterString(s []string, f func(string) bool) []string {
	var r[]string
	for _, v := range s {
		if f(v) {
			r = append(r, v)
		}
	}

	return r
}

func main() {
	evens := FilterInt([]int{1,2,3,4,5}, func(i int) bool { return i % 2 == 0 })
	fmt.Println("%v", evens)

	shortStrings := FilterString([]string{"ok", "notok", "maybe", "maybe not"}, func(s string) bool { return len(s) < 3 })
	fmt.Println("%v", shortStrings)
}

~/generics$
```

The example simply creates two lists and in turn filters them by a function we provide.

```
~/generics$ go run no_generics.go 

%v [2 4]
%v [ok]

~/generics$
```

While not the end of the world, our code is duplicated. Code duplication is the source of many copy-paste bugs. Ideally,
we would like to avoid those bugs. Worse, we might start thinking in terms of local optimizations just based on types. 
While this may benefit some types, in our example, probably no unique optimizations are needed.


### 2. Refactored with Generics

We now combine our two filter functions into a single function. Type parameters being used as constraints, a critical 
part of generics.

In our combined filter function, we abstract the type information with the constraint "[T any]". "any" is the same as
"interface{}" in prior editions of Go. The constraint tells the compiler which types are allowed and in turn generates
the required version of the code. 

You see in the parameter list we are limited to slice of type T and elements of type T are passed into a function we 
also provide, yielding a slice of the element type T.

```
~/generics$ vi with_generics.go

package main

import (
	"fmt"
)

func Filter[T any](s []T, f func(T) bool) []T {
	var r[]T
	for _, v := range s {
		if f(v) {
			r = append(r, v)
		}
	}

	return r
}

func main() {
	evens := Filter([]int{1,2,3,4,5}, func(i int) bool { return i % 2 == 0 })
	fmt.Println("%v", evens)

	shortStrings := Filter([]string{"ok", "notok", "maybe", "maybe not"}, func(s string) bool { return len(s) < 3 })
	fmt.Println("%v", shortStrings)
}

~/generics$
```

Running yields the same output as our prior version.

```
~/generics$ go run with_generics.go

%v [2 4]
%v [ok]

~/generics$
```

While we are sticking to examples, its worth thinking about the impact of generics on our code.

* line count  - 35 vs 24
* binary size - both around 1.8MB
* speed       - same (using time command)

Less code is typically considered a good thing. Potentially with any new feature, syntax can be challenging at first.


### 3. Custom Constraints

In this example, we look at a function (SoundOff) where the parameter is limited to the constraint Hybrid. The Hybrid
constraint ultimately is limited to the interface Animal. Our function SoundOff not only accepts all the Animals, it 
retreives the value type in order to do something specific (just to affirm we do not lose that ability). This isn't 
unique to generics, we have to do this with interfaces if you want access to the value type. Take note that we cast our
input parameter to the empty interface (any) in order to retrieve the concrete type. This casting is required as an 
interface is expected in the type switch, not a concrete type H.

```
~/generics$ vi simple_generics.go 

package main

import (
        "fmt"
)

type Animal interface {
        Sound()
}

type Cat struct{}

func (c Cat) Sound() { fmt.Println("Meow") }

func (c Cat) SpecialToCat() { fmt.Println("Cat special") }

type Dog struct{}

func (d Dog) Sound() { fmt.Println("Woof") }

func (c Dog) UniqueToDog() { fmt.Println("Dog unique") }

type Hybrid interface {
        Animal
}

func SoundOff[H Hybrid](animal H) H {
        animal.Sound()

        switch a := any(animal).(type) {
        case Dog:
                a.UniqueToDog()
        case Cat:
                a.SpecialToCat()
        default:
                fmt.Println("Then hoo?")
        }
        return animal
}

type NotAnOwl struct{} 

func (NotAnOwl) Sound() {}

func main() {
        var c Cat = SoundOff(Cat{})
        d := SoundOff(Dog{})

        c.Sound()
        c.SpecialToCat()
        d.Sound()
        d.UniqueToDog()

        SoundOff(NotAnOwl{})
}

~/generics$
```

Sit back and listen to our zoo.

```
~/generics$ go run simple_generics.go

Meow
Cat special
Woof
Dog unique
Meow
Cat special
Woof
Dog unique
Then hoo?

~/generics$
```

Per one of the lead developers, Ian Lance Taylor, here are guidelines on when to use generics.

When to use generics:

* Functions that work on slices, maps, channels of any element type, and have no assumptions about a particular element type is used
* General purpose data structures, ie. linked list, b-tree
 * Prefer functions versus methods (allows the data structure to remain agnostic to the type)
* When elements have a common method with the same implementation (Read(network) and Read(file) have different implementations so don't use generics)

When to not use generics:

* When just calling a method on the argument (use interfaces)
* When implementation of a common method differs
* When an operation differs per type (use reflection instead)


### 4. More on Collections

Collections such as arrays of "any" where we operate the same way on each item are prime targets for generics. In this 
example we use a channel and our own linked list to show how we can generify our iteration code. Unlike prior examples,
we have two constraints "MustBe" (our input) and "Result" (our output). Notice they do not have to match.

This example, we create a channel and hydrate it with some sample data; then launch a Go routine to iterate over it. 
Then we hand craft our own linked list, in turn, iterate over it. Notice our Iterate function takes either type (via the 
MustBe constraint) in turn extract a result.

```
~/generics$ cat iterate.go 

package main

import (
	"fmt"	
	"sync"
	"time"
)

// Linked List Example
type LL struct {
	N    *LL
	data string
}

func (l LL) Next() *LL { return l.N }

// Constraint using a channel and our Linked List type
type MustBe interface {
	chan string | LL
}

// Limit results to these types
type Result interface {
	string | LL
}

func Iterate[M MustBe, R Result](o M, iter func(M) R) (r R) {
	return iter(o)
}

func main() {
	c := make(chan string, 5)
	c <- "ok"
	c <- "ok2"

	citer := func(c chan string) string {
		select {
		case msg1 := <-c:
			return msg1
		case <-time.After(1 * time.Second):
			return "nothing"
		}
	}
	var wg sync.WaitGroup
	wg.Add(2)
	go func(f func(chan string) string) {
		for {
			fmt.Println(Iterate(c, f))
			wg.Done()
		}
	}(citer)

	wg.Wait()

	n1 := LL{data: "n1"}
	n2 := LL{data: "n2"}
	n3 := LL{data: "n3"}
	n1.N = &n2
	n2.N = &n3

	liter := func(l LL) LL {
		var zero LL

		if l.N != zero.N {
			return *l.N
		} else {
			return zero
		}
	}

	fmt.Println(Iterate(n1, liter))
	fmt.Println(Iterate(n2, liter))
	fmt.Println(Iterate(n3, liter))
}

~/generics$
```

When we launch, the channel is consumed first, then we work through the linked list.

```
~/generics$ go run iterate.go

[2 2 3 3]
ok
ok2
{0xc00000c0c0 n2}
{<nil> n3}
{<nil> }

~/generics$
```

With generics, we are able to use a single function to handle multiple input types. The use of a function to retrieve 
the next element means we don't have to support a common interface with these types. This means in theory we can use
other types as well.


### 5. Composed Type

We can begin by trying to use Generics to remove the triplicated check, via the following constraint.

```
type MyType interface {
    bytes.Buffer | bytes.Reader | strings.Reader
}
```

We then change the signature to the following:

`func NewRequest[M MyType](method, target string, body M) *http.Request {`

This is a start, but we are missing the following:

* We are unable to simply check for nil (before the switch we see "if body != nil") with generics
* The compiler will tell us Len() is missing (since its not part of our compostition of types in our constraint)

We can address the first issue with checking a related zero value instead of nil. The original nil check becomes the 
following. Remember, generics will identify the arguments (and the actual types) to fill in what M should be.

```
...
	var zero M
	if body != zero {
...        
```

This still leaves us with the case of the missing Len(). Originally, the parameter was an interface, "io.Reader", which
doesn't have a Len() method. This drove the original need to use a switch statement to access the value types, which 
did each implemented a Len() method. We could for example make a new interface, one that combines Reader and Len().

```
type LenReader interface {
        io.Reader
        Len() int
}
```

This however doesn't do the job, we are still exposed to any type that implements that interface (not limited at compiler time).

In the following code, we take the NewRequest fuction (borrowed from httptest.go) and refactor it.

```
~/generics$ vi complex_generics.go

package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type myStruct struct {
	s *strings.Reader
}

func (m myStruct) Len() int {
	return m.s.Len()
}

func (m myStruct) Read(b []byte) (int, error) {
	return m.s.Read(b)
}

type evil struct { }
func (e evil) Len() int { return -1 }
func (e evil) Read(b []byte) (int, error) {
	return -1, nil
}

type MyType interface {
	*bytes.Buffer | *bytes.Reader | *strings.Reader | myStruct // | evil 

	Len() int
	io.Reader
	comparable
}

type Lener interface {
	Len() int
}

// ./http/httptest/httptest.go
func NewRequest[M MyType](method, target string, body M) *http.Request {
	if method == "" {
		method = "GET"
	}
	req, err := http.ReadRequest(bufio.NewReader(strings.NewReader(method + " " + target + " HTTP/1.0\r\n\r\n")))
	if err != nil {
		panic("invalid NewRequest arguments; " + err.Error())
	}

	// HTTP/1.0 was used above to avoid needing a Host field. Change it to 1.1 here.
	req.Proto = "HTTP/1.1"
	req.ProtoMinor = 1
	req.Close = false

	var zero M
	if body != zero {
		switch i := any(body).(type) {
		case Lener, io.ReadCloser:
			if b, ok := i.(Lener); ok {
				req.ContentLength = int64(b.Len())
			}
			if rc, ok := i.(io.ReadCloser); ok {
				req.Body = rc
			}
		default:
			req.Body = io.NopCloser(body)
		}
	} else {
		req.ContentLength = -1
	}

	// 192.0.2.0/24 is "TEST-NET" in RFC 5737 for use solely in
	// documentation and example source code and should not be
	// used publicly.
	req.RemoteAddr = "192.0.2.1:1234"

	if req.Host == "" {
		req.Host = "example.com"
	}

	if strings.HasPrefix(target, "https://") {
		req.TLS = &tls.ConnectionState{
			Version:           tls.VersionTLS12,
			HandshakeComplete: true,
			ServerName:        req.Host,
		}
	}

	return req
}

func main() {
	fmt.Println(NewRequest("GET", "/", myStruct{strings.NewReader("")}))
	fmt.Println(NewRequest("GET", "/", myStruct{}))
	fmt.Println(NewRequest("GET", "/", strings.NewReader("")))
	fmt.Println(NewRequest("GET", "/", &bytes.Buffer{}))
	fmt.Println(NewRequest("GET", "/", bytes.NewReader([]byte("read me"))))
//	fmt.Println(NewRequest("GET", "/", evil{}))
}

~/generics$
```

Running the example works (output not the important part here).

```
~/generics$ go run complex_generics.go 

&{GET / HTTP/1.1 1 1 map[] {} <nil> 0 [] false example.com map[] map[] <nil> map[] 192.0.2.1:1234 / <nil> <nil> <nil> <nil>}
...

~/generics$
```

By moving our types and method requirements to our MyType constraint we are now limited to a subset of types versus all
types that implement io.Reader. 

I think most would argue this code is less readable than before. Should generics be used in this case? Based on other 
examples, it most likely should remain using io.Reader and a type switch. One could use the following alternative to 
manage Len calls.

```
...
	if body != nil {
		if b, ok := body.(interface{ Len() int }); ok {
			req.ContentLength = int64(b.Len())
		}
		if rc, ok := body.(io.ReadCloser); ok {
			req.Body = rc
		} else {
			req.Body = io.NopCloser(body)
		}
	} else {
		req.ContentLength = -1
	}
...	
```


### 6. Conclusion

Generics are here, and will only improve. What we have in Go 1.18 is not even the final take, new generic functions are to be included in the standard library along with potential other improvements. Until then keep on eye on the official blog(s) such as this one:

https://go.dev/blog/why-generics

<br>

Congratulations you have completed the lab!!

<br>

