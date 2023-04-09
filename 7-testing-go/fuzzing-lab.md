
# Go


## Fuzzing

Fuzzing is a technique where one automatically generates inputs, which in turn causes the receiving program to respond
in some unexpected way. How a fuzzing program generates inputs is often classified as either simple to complex spectrum.
For example the inputs could simply be randomily generated, or you can provide an initial set of values, or even more
complex fuzzing systems can be rule driven (including deriving new rules that generate potential values).


### 1. Fuzzing in Go

Third party fuzzing software has been around (even for Go), but in Go 1.18 the technique was integrated into the testing
package.

The goal for us is to learn the basic mechanics of using Go's fuzzing functionality. Before we look at code, take some
time to review the inital documentation.

`~$ go doc testing.F`

Our first example leverages the existing documentation example. As you read through this example:

* Note how testing.F is used along with testing.T
* Note the use of seed data (more on this later)
* Think about the actual test, is it a regular looking unit test (yes!)

Create the first example.

```
~$ cd

~$ mkdir ~/fuzzing

~$ cd appendix

~/fuzzing$ vi fuzzing_test.go

package main

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func FuzzMe(f *testing.F) {
	for _, seed := range [][]byte{{}, {0}, {9}, {0xa}, {0xf}, {1, 2, 3, 4}} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, in []byte) {
		enc := hex.EncodeToString(in)
		out, err := hex.DecodeString(enc)
		if err != nil {
			t.Fatalf("%v: decode: %v", in, err)
		}
		if !bytes.Equal(in, out) {
			t.Fatalf("%v: not equal after round trip: %v", in, out)
		}
	})
}

~/fuzzing$
```

To run the test using the embedded seed corpus.

```
~/fuzzing$ go test -v --test.fuzz "" fuzzing_test.go

=== RUN   FuzzMe
=== RUN   FuzzMe/seed#0
=== RUN   FuzzMe/seed#1
=== RUN   FuzzMe/seed#2
=== RUN   FuzzMe/seed#3
=== RUN   FuzzMe/seed#4
=== RUN   FuzzMe/seed#5
--- PASS: FuzzMe (0.00s)
    --- PASS: FuzzMe/seed#0 (0.00s)
    --- PASS: FuzzMe/seed#1 (0.00s)
    --- PASS: FuzzMe/seed#2 (0.00s)
    --- PASS: FuzzMe/seed#3 (0.00s)
    --- PASS: FuzzMe/seed#4 (0.00s)
    --- PASS: FuzzMe/seed#5 (0.00s)
PASS
ok  	command-line-arguments	0.105s

~/fuzzing$
```

As you can see, we have six runs executed, this maps to the number of seeds included via `f.Add(seed)`. By including 
seeds we are afforded the opportunity to have regression protection, just like a typical unit test. This mode of 
execution happened because we did not specify the test to fuzz. An alternative way to execute the same thing (using
seeded versus generated data) is:

`go test -v --test.run "FuzzMe" fuzzing_test.go`

> nb. Notice its `--test.run` and not `--test.fuzz`.

For example, using seeded data allows the test to perform fast enough to be part of a CI pipeline.

While useful, we haven't actually fuzzed anything, we simply ran our basic comparisons. To fuzz, we need to specify our
tests we wanted fuzzed. To generate new results (aka fuzz), specify the test (grab a drink, its going to be a long time
-- if ever -- to complete). When you are tired of waiting, hit control c.

```
~/fuzzing$ go test -v --test.fuzz "FuzzMe" fuzzing_test.go

=== FUZZ  FuzzMe
fuzz: elapsed: 0s, gathering baseline coverage: 0/29 completed
fuzz: elapsed: 0s, gathering baseline coverage: 29/29 completed, now fuzzing with 12 workers
fuzz: elapsed: 3s, execs: 780426 (260128/sec), new interesting: 0 (total: 29)
fuzz: elapsed: 6s, execs: 1577625 (265663/sec), new interesting: 0 (total: 29)
fuzz: elapsed: 9s, execs: 2369952 (264185/sec), new interesting: 0 (total: 29)
fuzz: elapsed: 12s, execs: 3024147 (218055/sec), new interesting: 0 (total: 29)
fuzz: elapsed: 15s, execs: 3759791 (245166/sec), new interesting: 0 (total: 29)
fuzz: elapsed: 18s, execs: 4470924 (237062/sec), new interesting: 0 (total: 29)
fuzz: elapsed: 21s, execs: 5178109 (235752/sec), new interesting: 0 (total: 29)
fuzz: elapsed: 24s, execs: 5898784 (240167/sec), new interesting: 0 (total: 29)
^Cfuzz: elapsed: 25s, execs: 6211453 (222820/sec), new interesting: 0 (total: 29)
--- PASS: FuzzMe (25.41s)
PASS
ok  	command-line-arguments	25.545s

~/fuzzing$
```

The immediate change is we specifiy a test to fuzz via `--test.fuzz "FuzzMe"`. Instead of just running the seeds, we are 
now generating new inputs and testing them. This process will continue for a very long time, hence the need to stop it.
The reason it takes (hours, days, more!) is because the range of inputs can be huge (just a single integer has billions
of potenial values). This computational expense is why fuzzing is not happening by default.

We do have some additional controls though, versus waiting forever:

* Fuzzing stops when a test fails
* Limit the time to fuzz via -fuzztime (eg. -fuzztime 30s)
* Hit control^c or send SIGINT signal (eg. in a CI pipeline)


### 2. Fuzzing Failure

In our prior example, we did not see a failure, however, we ran for only under a minute, it could take days to find bad
input. 

Lets now create a bad test to highlight the process, and what we do with fuzzing generated test inputs after finding 
them! The following test will error if the value is greater than 100 or less than 1000. Normally, the invalid inputs 
are not so obvious. Add to your existing test file.

```
~/fuzzing$ vi fuzzing_test.go 

... existing code ...

func FuzzBad(f *testing.F) {
	f.Fuzz(func(t *testing.T, i int) {
		if i > 100 && i < 1000 {
			t.Fatalf("want: 101-999, got: %v", i)
		}	
	})
}

~/fuzzing$
```

Running the fuzzing directly against this test produces the following.

```
~/fuzzing$ go test -v --test.fuzz "FuzzBad" fuzzing_test.go

=== RUN   FuzzMe
=== RUN   FuzzMe/seed#0
=== RUN   FuzzMe/seed#1
=== RUN   FuzzMe/seed#2
=== RUN   FuzzMe/seed#3
=== RUN   FuzzMe/seed#4
=== RUN   FuzzMe/seed#5
--- PASS: FuzzMe (0.00s)
    --- PASS: FuzzMe/seed#0 (0.00s)
    --- PASS: FuzzMe/seed#1 (0.00s)
    --- PASS: FuzzMe/seed#2 (0.00s)
    --- PASS: FuzzMe/seed#3 (0.00s)
    --- PASS: FuzzMe/seed#4 (0.00s)
    --- PASS: FuzzMe/seed#5 (0.00s)
=== FUZZ  FuzzBad
fuzz: elapsed: 0s, gathering baseline coverage: 0/1 completed
fuzz: elapsed: 0s, gathering baseline coverage: 1/1 completed, now fuzzing with 12 workers
fuzz: elapsed: 0s, execs: 9 (530/sec), new interesting: 0 (total: 1)
--- FAIL: FuzzBad (0.02s)
    --- FAIL: FuzzBad (0.00s)
        fuzzing_test.go:28: want: 101-999, got: 174
    
    Failing input written to testdata/fuzz/FuzzBad/daef2fa4fc63690477c788772f4488eb55a67946e4bed16916b63688c2c99935
    To re-run:
    go test -run=FuzzBad/daef2fa4fc63690477c788772f4488eb55a67946e4bed16916b63688c2c99935
FAIL
exit status 1
FAIL	command-line-arguments	0.151s

~/fuzzing$
```

We see the original test still runs over the provided seeds. This is by design, meant to run the seeds as a regression
protective measure. If no seeds where available, the test would not run without using direct invocation via --test.fuzz.

More importantly, we see a failing set of input(s) was identified. In the example, the value 174 was attempted and 
failed. Yours may differ. The result is stored in a directory called testdata. This directory will in turn be used as
the 'seed' for future regression runs.

We see the directory structure as follows.

```
~/fuzzing$ tree ./testdata

testdata
└── fuzz
    └── FuzzBad
        └── daef2fa4fc63690477c788772f4488eb55a67946e4bed16916b63688c2c99935

2 directories, 1 file

~/fuzzing$
```

Reviewing the auto-generated file, we see the encoded bad input.

```
~/fuzzing$ cat ./testdata/fuzz/FuzzBad/daef2fa4fc63690477c788772f4488eb55a67946e4bed16916b63688c2c99935 

go test fuzz v1
int(174)

~/fuzzing$
```

In our result, it is simply an int, but for more complex results we will see more complex encodings.

If you run the test again, it simply fails again, treating the ./testdata now as a seed.

```
~/fuzzing$ go test --test.fuzz "FuzzBad" fuzzing_test.go  

fuzz: elapsed: 0s, gathering baseline coverage: 0/2 completed
failure while testing seed corpus entry: FuzzBad/daef2fa4fc63690477c788772f4488eb55a67946e4bed16916b63688c2c99935
fuzz: elapsed: 0s, gathering baseline coverage: 0/2 completed
--- FAIL: FuzzBad (0.02s)
    --- FAIL: FuzzBad (0.00s)
        fuzzing_test.go:28: want: 101-999, got: 174
    
FAIL
exit status 1
FAIL	command-line-arguments	0.126s
```

Lets almost fix this, lets simply reduce from 1000 to 150, shrinking the potential error range. If your error was less
than 150, change the range to allow it.

```
vi fuzzing_test.go 

... update FuzzBad with below ...

func FuzzBad(f *testing.F) {
        f.Fuzz(func(t *testing.T, i int) {
                if i > 100 && i < 150 {
                        t.Fatalf("want: 101-150, got: %v", i)
                }
        })
}

~/fuzzing$
```

Running once more.

```
~/fuzzing$ go test --test.fuzz "FuzzBad" fuzzing_test.go

fuzz: elapsed: 0s, gathering baseline coverage: 0/2 completed
fuzz: elapsed: 0s, gathering baseline coverage: 2/2 completed, now fuzzing with 12 workers
fuzz: elapsed: 0s, execs: 4 (205/sec), new interesting: 0 (total: 2)
--- FAIL: FuzzBad (0.02s)
    --- FAIL: FuzzBad (0.00s)
        fuzzing_test.go:28: want: 101-150, got: 117
    
    Failing input written to testdata/fuzz/FuzzBad/7bd545fe2a8997effdf791253ba576f785189c92f46c205024dc835aa7f63b27
    To re-run:
    go test -run=FuzzBad/7bd545fe2a8997effdf791253ba576f785189c92f46c205024dc835aa7f63b27
FAIL
exit status 1
FAIL	command-line-arguments	0.234s
```

Reviewing the ./testdata, we see the new failure.

```
~/fuzzing$ tree testdata

testdata
└── fuzz
    └── FuzzBad
        ├── 7bd545fe2a8997effdf791253ba576f785189c92f46c205024dc835aa7f63b27
        └── daef2fa4fc63690477c788772f4488eb55a67946e4bed16916b63688c2c99935

2 directories, 2 files

~/fuzzing$ cat testdata/fuzz/FuzzBad/7bd545fe2a8997effdf791253ba576f785189c92f46c205024dc835aa7f63b27 

go test fuzz v1
int(117)

~/fuzzing$
```

As expected, a new failing case was identified and stored. If we removed the failing logic, can we fuzz?

```
~/fuzzing$ vi fuzzing_test.go

...

func FuzzBad(f *testing.F) {
        f.Fuzz(func(t *testing.T, i int) {
                # hope this never fails
                if i != i {
                        t.Fatalf("want: %v, got: %v", i, i)
                }
        })
}

~/fuzzing$
```

Running with a 10 second fuzzing time limit.

```
~/fuzzing$ go test --test.fuzz "FuzzBad" fuzzing_test.go --test.fuzztime 10s

fuzz: elapsed: 0s, gathering baseline coverage: 0/3 completed
fuzz: elapsed: 0s, gathering baseline coverage: 3/3 completed, now fuzzing with 12 workers
fuzz: elapsed: 3s, execs: 794663 (264831/sec), new interesting: 0 (total: 3)
fuzz: elapsed: 6s, execs: 1608866 (271422/sec), new interesting: 0 (total: 3)
fuzz: elapsed: 9s, execs: 2400547 (263916/sec), new interesting: 0 (total: 3)
fuzz: elapsed: 10s, execs: 2642809 (217213/sec), new interesting: 0 (total: 3)
PASS
ok  	command-line-arguments	10.273s

~/fuzzing$
```

We still see the seeds (from previous captured results) are executed; and we fuzz for an additional 10 seconds.


### 3. Conclusion

While the fuzzing technique is not new, its arrival in Go 1.18 means it will be some time before we see more intersting
examples. More about fuzzing in Go can be learned from here https://go.dev/doc/tutorial/fuzz.

Fuzzing is not enough to prove program correctness, for that formal methods must be used. Features of Go itself are
derived from work relating to formal methods. Channels (a key feature of Go) are derived from communicating sequential processes (CSP). CSP is a formal language describing concurrent systems (eg. Using channels between Go routines).

<br>

Congratulations you have completed the lab!!

<br>
