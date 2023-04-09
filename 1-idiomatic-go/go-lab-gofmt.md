
# Go


## Diving Deep with Gofmt

Each Go command or subcommand has various ways to view and how to execute or use them, examples including `go fmt -h`,
`go help fmt`, `go doc cmd/gofmt`.  Some commands/subcommands have wrappers/aliases (ex. `go fmt` calls `gofmt`), the
wrapper may include its own flags and behavior.


## Preperation

If not installed, install the latest version of Go.

```
~$ curl -sLO https://go.dev/dl/go1.20.3.linux-amd64.tar.gz

~$ sudo tar -C /usr/local -xzf go1.20.3.linux-amd64.tar.gz

~$ export PATH=/usr/local/go/bin:$PATH

~$ go version

go version go1.20.3 linux/amd64

~$
```


## gofmt

Gofmt formats Go programs. It uses tabs for indentation and blanks for alignment. Alignment assumes that an editor is
using a fixed-width font.

Beyond basic formatting, gofmt includes features to simplify some complex coding syntax (-s) including user defined
replacement rules (-r rule).

Gofmt takes arguments to simply notify if a modification would take place (-l) and too display deltas (-d).  See the
output of `go doc cmd/gofmt` for details.

In the following examples, we reference the actual test data used by `gofmt`.  When Go is installed (ex.
`/usr/local/go/`), it will include the standard library code and related testing code/data.

The `gofmt` source and test data can be found in Go installation, ex. `$(go env GOROOT)/src/cmd/gofmt`.

```
~$ cd ~

~$ ls -l /usr/local/go/src/cmd/gofmt

total 64
-rw-r--r-- 1 root root  3239 May 10 16:48 doc.go
-rw-r--r-- 1 root root 14942 May 10 16:48 gofmt.go
-rw-r--r-- 1 root root  6144 May 10 16:48 gofmt_test.go
-rw-r--r-- 1 root root  5166 May 10 16:48 internal.go
-rw-r--r-- 1 root root  4011 May 10 16:48 long_test.go
-rw-r--r-- 1 root root  8346 May 10 16:48 rewrite.go
-rw-r--r-- 1 root root  4774 May 10 16:48 simplify.go
drwxr-xr-x 2 root root  4096 May 10 16:48 testdata

~$
```

Take a moment to look under the `.../testdata/` directory.  

* files ending in `input` are the files used as test input
* files ending in `golden` are the results to compare against after running `gofmt`

Feel free to review `gofmt_test.go` on how it uses this data.

```
~$ ls -l $(go env GOROOT)/src/cmd/gofmt/testdata

total 232
-rw-r--r-- 1 root root   76 May 10 16:48 comments.golden
-rw-r--r-- 1 root root   76 May 10 16:48 comments.input
-rw-r--r-- 1 root root 2596 May 10 16:48 composites.golden
-rw-r--r-- 1 root root 3236 May 10 16:48 composites.input
-rw-r--r-- 1 root root  234 May 10 16:48 crlf.golden
-rw-r--r-- 1 root root  247 May 10 16:48 crlf.input
-rw-r--r-- 1 root root  140 May 10 16:48 emptydecl.golden
-rw-r--r-- 1 root root  148 May 10 16:48 emptydecl.input
-rw-r--r-- 1 root root 1945 May 10 16:48 go2numbers.golden
-rw-r--r-- 1 root root 2004 May 10 16:48 go2numbers.input
-rw-r--r-- 1 root root 2147 May 10 16:48 import.golden
-rw-r--r-- 1 root root 2152 May 10 16:48 import.input
-rw-r--r-- 1 root root  249 May 10 16:48 issue28082.golden
-rw-r--r-- 1 root root  447 May 10 16:48 issue28082.input
-rw-r--r-- 1 root root  307 May 10 16:48 ranges.golden
-rw-r--r-- 1 root root  304 May 10 16:48 ranges.input
-rw-r--r-- 1 root root  248 May 10 16:48 rewrite1.golden
-rw-r--r-- 1 root root  248 May 10 16:48 rewrite1.input
...
-rw-r--r-- 1 root root  228 May 10 16:48 rewrite9.golden
-rw-r--r-- 1 root root  238 May 10 16:48 rewrite9.input
-rw-r--r-- 1 root root  927 May 10 16:48 slices1.golden
-rw-r--r-- 1 root root  957 May 10 16:48 slices1.input
-rw-r--r-- 1 root root   32 May 10 16:48 stdin1.golden
-rw-r--r-- 1 root root   32 May 10 16:48 stdin1.input
...
-rw-r--r-- 1 root root  342 May 10 16:48 stdin7.golden
-rw-r--r-- 1 root root  296 May 10 16:48 stdin7.input
-rw-r--r-- 1 root root  365 May 10 16:48 typealias.golden
-rw-r--r-- 1 root root  360 May 10 16:48 typealias.input
-rw-r--r-- 1 root root  680 May 10 16:48 typeparams.golden
-rw-r--r-- 1 root root  677 May 10 16:48 typeparams.input
-rw-r--r-- 1 root root 1476 May 10 16:48 typeswitch.golden
-rw-r--r-- 1 root root 1473 May 10 16:48 typeswitch.input

~$
```

As you can see, there are tests for hidden characters (ex. newlines on different operating systems), to various
structural rewrite rules, to syntax switching.  We will first compare `gofmt` and `go fmt` while looking at the test for
line endings.  Because line endings are not typically visable, we will leverage `cat -t` to help view the differences.

```
~$ cat -t $(go env GOROOT)/src/cmd/gofmt/testdata/crlf.input

/*^M
^ISource containing CR/LF line endings.^M
^IThe gofmt'ed output must only have LF^M
^Iline endings.^M
^ITest case for issue 3961.^M
*/^M
package main^M
^M
func main() {^M
^I// line comment^M
^Iprintln("hello, world!") // another line comment^M
^Iprintln()^M
}^M

~$
```

Some test input actually describes or links to details about the issue being fixed (via formatting).  Lets use `gofmt`
to see the potential changes (`-d` shows the pending delta).

```
~$ gofmt -d $(go env GOROOT)/src/cmd/gofmt/testdata/crlf.input | cat -t -

diff -u /usr/local/go/src/cmd/gofmt/testdata/crlf.input.orig /usr/local/go/src/cmd/gofmt/testdata/crlf.input
--- /usr/local/go/src/cmd/gofmt/testdata/crlf.input.orig^I2022-06-01 02:43:33.809159190 +0000
+++ /usr/local/go/src/cmd/gofmt/testdata/crlf.input^I2022-06-01 02:43:33.809159190 +0000
@@ -1,13 +1,13 @@
-/*^M
-^ISource containing CR/LF line endings.^M
-^IThe gofmt'ed output must only have LF^M
-^Iline endings.^M
-^ITest case for issue 3961.^M
-*/^M
-package main^M
-^M
-func main() {^M
-^I// line comment^M
-^Iprintln("hello, world!") // another line comment^M
-^Iprintln()^M
-}^M
+/*
+^ISource containing CR/LF line endings.
+^IThe gofmt'ed output must only have LF
+^Iline endings.
+^ITest case for issue 3961.
+*/
+package main
+
+func main() {
+^I// line comment
+^Iprintln("hello, world!") // another line comment
+^Iprintln()
+}

~$
```

By using `cat -t` we can see the difference.  Many CI pipelines use `gofmt -d ...` to check changes meet project
standards.

How would this work/look when running through `go fmt`.

```
~$ go fmt $(go env GOROOT)/src/cmd/gofmt/testdata/crlf.input

cannot find package "." in:
	/usr/local/go/src/cmd/gofmt/testdata/crlf.input

~$
```

That is unfortunate, what happened?  As mentioned before, wrappers may have their own logic.  Review the short help for
`go fmt`.  

```
~$ go fmt -h

usage: go fmt [-n] [-x] [packages]
Run 'go help fmt' for details.

~$
```

In this case, it appears the `go fmt` doesn't want to just load a file, instead it will do some basic sematic analysis
and load only proper packages.  We can investigate that hunch by reviewing the code.

```
~$ grep -i package $(go env GOROOT)/src/cmd/go/internal/fmtcmd/fmt.go

// Package fmtcmd implements the ``go fmt'' command.
package fmtcmd
	UsageLine: "go fmt [-n] [-x] [packages]",
	Short:     "gofmt (reformat) package sources",
Fmt runs the command 'gofmt -l -w' on the packages named
For more about specifying packages, see 'go help packages'.
	for _, pkg := range load.PackagesAndErrors(ctx, load.PackageOpts{}, args) {
				fmt.Fprintf(os.Stderr, "go: not formatting packages in dependency modules\n")
		// the command only applies to this package,
		// not to packages in subdirectories.

~$
```

If time allows, take a look at `$(go env GOROOT)/src/cmd/go/internal/fmtcmd/fmt.go`. For now, we will trust it is
blocking our initial attempt. Lets fix this by moving the sample input into a package.

```
~$ grep package $(go env GOROOT)/src/cmd/gofmt/testdata/crlf.input

package main

~$
```

Since this code is in package "main", we just need to make a package directory that doesn't conflict with any existing
in our GOPATH.

```
~$ TEMPPKGDIR=main$RANDOM && echo $TEMPPKGDIR

main30393

~$ mkdir -p $(go env GOPATH)/src/$TEMPPKGDIR

~$ cp $(go env GOROOT)/src/cmd/gofmt/testdata/crlf.input $(go env GOPATH)/src/$TEMPPKGDIR/crlf.go

~$ go fmt -x $(go env GOPATH)/src/$TEMPPKGDIR/crlf.go

/usr/local/go/bin/gofmt -l -w go/src/main30393/crlf.go
go/src/main30393/crlf.go

~$
```

There is no exciting output as there is not much in the file!

Notice we changed not only the location, but the file name.  `go fmt` is stricter when loading files then `gofmt` (aka `.input` -> `.go`).

To run all the tests related to `gofmt`, we can as follows.

```
~$ go test --count=1 -v $(go env GOROOT)/src/cmd/gofmt/...

=== RUN   TestRewrite
--- PASS: TestRewrite (0.02s)
=== RUN   TestCRLF
--- PASS: TestCRLF (0.00s)
=== RUN   TestBackupFile
    gofmt_test.go:195: Created: /tmp/gofmt_test958648301/foo.go1653915020
--- PASS: TestBackupFile (0.00s)
=== RUN   TestDiff
--- PASS: TestDiff (0.00s)
=== RUN   TestReplaceTempFilename
--- PASS: TestReplaceTempFilename (0.00s)
=== RUN   TestAll
    long_test.go:92: known gofmt idempotency bug (Issue #24472)
--- PASS: TestAll (4.77s)
PASS
ok  	cmd/gofmt	4.801s

~$
```

We are using `--count=1` to disable caching (this is the official way verus GOCACHE=off).

To learn more about `gofmt`, lets see what else is in the testdata directory.  For example, review `ranges.input`.

```
~$ cat $(go env GOROOT)/src/cmd/gofmt/testdata/ranges.input

//gofmt -s

// Test cases for range simplification.
package p

func _() {
	for a, b = range x {}
	for a, _ = range x {}
	for _, b = range x {}
	for _, _ = range x {}

	for a = range x {}
	for _ = range x {}

	for a, b := range x {}
	for a, _ := range x {}
	for _, b := range x {}

	for a := range x {}
}

~$
```

What's wrong with that code, doesn't look so bad?  Lets run it.

```
~$ go run $(go env GOROOT)/src/cmd/gofmt/testdata/ranges.input

cannot find package "." in:
	/usr/local/go/src/cmd/gofmt/testdata/ranges.input

~$
```

We can see from the prior listing that it was never going to work.  First, there is no `main` function, second, all the
variables appear ill-defined.  For now, lets see if `gofmt` will still apply some magic.  Using the flags `-d` (show
diff) and `-s` to apply simplification rules, lets check the results.

```
~$ gofmt -d -s $(go env GOROOT)/src/cmd/gofmt/testdata/ranges.input

diff -u /usr/local/go/src/cmd/gofmt/testdata/ranges.input.orig /usr/local/go/src/cmd/gofmt/testdata/ranges.input
--- /usr/local/go/src/cmd/gofmt/testdata/ranges.input.orig	2022-06-01 02:47:00.416235165 +0000
+++ /usr/local/go/src/cmd/gofmt/testdata/ranges.input	2022-06-01 02:47:00.416235165 +0000
@@ -4,17 +4,27 @@
 package p
 
 func _() {
-	for a, b = range x {}
-	for a, _ = range x {}
-	for _, b = range x {}
-	for _, _ = range x {}
+	for a, b = range x {
+	}
+	for a = range x {
+	}
+	for _, b = range x {
+	}
+	for range x {
+	}
 
-	for a = range x {}
-	for _ = range x {}
+	for a = range x {
+	}
+	for range x {
+	}
 
-	for a, b := range x {}
-	for a, _ := range x {}
-	for _, b := range x {}
+	for a, b := range x {
+	}
+	for a := range x {
+	}
+	for _, b := range x {
+	}
 
-	for a := range x {}
+	for a := range x {
+	}
 }

~$
```

It still appears to be working, even though the file is not a `go` extension nor in a proper directory.  In this case,
we can see a transformation, so what happened?  The very first line in the `input` file is instructing `gofmt` test
suite to use the simplify flag.

```
~$ head $(go env GOROOT)/src/cmd/gofmt/testdata/ranges.input

//gofmt -s

// Test cases for range simplification.
package p

func _() {
	for a, b = range x {}
	for a, _ = range x {}
	for _, b = range x {}
	for _, _ = range x {}

~$
```

We explicitly called simplify (`-s`) since we are not using the test harness to run this input.

The `-s` flag performs the following (see `go doc cmd/gofmt`):

```
...

    -s
    	Try to simplify code (after applying the rewrite rule, if any).

...

The simplify command

When invoked with -s gofmt will make the following source transformations
where possible.

    An array, slice, or map composite literal of the form:
    	[]T{T{}, T{}}
    will be simplified to:
    	[]T{{}, {}}

    A slice expression of the form:
    	s[a:len(s)]
    will be simplified to:
    	s[a:]

    A range of the form:
    	for x, _ = range v {...}
    will be simplified to:
    	for x = range v {...}

    A range of the form:
    	for _ = range v {...}
    will be simplified to:
    	for range v {...}

This may result in changes that are incompatible with earlier versions of
Go.

~$
```

It is a little challenging to view the diff vertically so lets view the results horizontally.

```
~$ diff -y -w \
<(cat $(go env GOROOT)/src/cmd/gofmt/testdata/ranges.input) \
<(gofmt -s $(go env GOROOT)/src/cmd/gofmt/testdata/ranges.input)

//gofmt -s							//gofmt -s

// Test cases for range simplification.				// Test cases for range simplification.
package p							package p

func _() {							func _() {
	for a, b = range x {}				      |		for a, b = range x {
	for a, _ = range x {}				      |		}
	for _, b = range x {}				      |		for a = range x {
	for _, _ = range x {}				      |		}
							      >		for _, b = range x {
							      >		}
							      >		for range x {
							      >		}

	for a = range x {}				      |		for a = range x {
	for _ = range x {}				      |		}
							      >		for range x {
							      >		}

	for a, b := range x {}				      |		for a, b := range x {
	for a, _ := range x {}				      |		}
	for _, b := range x {}				      |		for a := range x {
							      >		}
							      >		for _, b := range x {
							      >		}

	for a := range x {}				      |		for a := range x {
							      >		}
}								}

~$
```

We are using the `diff` command which has the ability to show a diff side by side (via `-y` flag).  We also use `-w` to
ignore any white space differences.  If you review the `gofmt` code, you will also see they use the `diff` tool when
showing `-d` output.  `gofmt` also relies on `diff`.

Based on the description of the simplify rules, it appears to be working, even if the semantic parsing was not applied.

> To learn a little more on gofmt and related tooling see the following related isses:
> https://github.com/golang/go/issues/27099
> https://github.com/golang/go/issues/27166
> https://go-review.googlesource.com/c/go/+/153459/


### Challenge

Review other `.input` files and see which flags are applied for example during rewrite tests.  Try to execute these
other inputs againt `gofmt` or `go fmt` to see how they work.


<br>

Congratulations you have completed the lab!

<br>

