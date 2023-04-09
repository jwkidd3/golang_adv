
# Go


## Context

The Go standard library includes a context package designed to allow a caller to pass along call context (metadata) and
to cancel or enforce a timeout on potentially long running call sequences.

Go microservices generally handle each incoming request using a separate goroutine. Request handlers often start
additional goroutines to perform high latency operations, such as accessing backends like databases and RPC services.
For example, rather than waiting for a DB to return before writing a message to Kafka, we can start two go routines, one
for the DB operation and one for the Kafka operation. This will allow the Kafka and DB activity to occur concurrently,
in some cases cutting the wait time in half.

The set of goroutines working on a request may need access to request-specific data such as the identity of the end
user, authorization tokens, and the request's deadline. When a request is canceled or times out, all the goroutines
working on that request should exit quickly so the system can reclaim resources and return to the end user promptly.
This is the exact use case Go context was designed for.

The Go context package Context type allows context metadata to be safely propagated through a graph of async and
concurrent goroutines. The context.Context also supports operation cancelation and deadlines. As a part of the Go
standard library, a wide range of built-in and third party packages/frameworks can easily support context operations
consistently.

Once a Go context is established it should be propagated to all subordinate goroutines so that cancelations can be
handled by all associated concurrent blocks of code. While Go context only works within a single process, it models the
behavior of HTTP headers in many ways (and with a little work by the developer, can be extended to operate over a
network API).

In this lab we will create a simple http service with multiple goroutines that uses Go context to propagate call
metadata and manage early termination. Toward the end of the lab we'll add support for propagating HTTP headers to all
of our goroutines.



### 1. Create a simple http service

To experiment with the context package let's create a simple http based service that responds to a caller with the route
hit. Enter the following program in a source file called ctx.go, run it:

```
~$ mkdir ~/context && cd $_

~/context$ vi ctx.go

package main

import (
        "fmt"
        "log"
        "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hit on %s!\n", r.URL.Path[1:])
}

func main() {
        http.HandleFunc("/", handler)
        log.Fatal(http.ListenAndServe(":4444", nil))
}

~/context$ go run ctx.go

```

In a second terminal, test the service:

```
~$ curl localhost:4444/test

Hit on test!

~$ curl localhost:4444/catamaran

Hit on catamaran!

~$
```

Great! We have a basic service up and running. Coming up, we'll make the service a little more complex and add context
support.


### 2. Add an internal goroutine

Imagine our service from the previous step would like to store data in a cloud database. Imagine this cloud database is
a best efforts system and can require an extended period of time to store data when busy. We want to perform the
database write asynchronously, returning to the user if the db write takes too long.

A first cut at a solution might be to use a goroutine to write to the slow database. That way our main service handler
can return to the caller at will. Add a mock dbwriter() goroutine to your service, updating your code as follows:

```go
~/context$ vi ctx.go

package main

import (
        "fmt"
        "log"
        "math/rand"
        "net/http"
        "strconv"
        "time"
)

func dbwriter(route string, done chan bool) {
        log.Println("db save " + route + " START")
        time.Sleep(time.Duration(rand.Intn(5)) * time.Second) //simulate db write by sleeping 0-4 seconds
        log.Println("db save " + route + " FINISH")
        done <- true
}

func handler(w http.ResponseWriter, r *http.Request) {
        route := r.URL.Path[1:]
        done := make(chan bool)
        go dbwriter(route, done)
        select {
        case <-done:
                w.WriteHeader(http.StatusCreated)
                fmt.Fprintf(w, "Hit on %s!\n", route)
        case <-time.After(2 * time.Second):
                log.Println("db timeout")
                w.WriteHeader(http.StatusGatewayTimeout)
        }
}

func main() {
        port := 4444
        rand.Seed(time.Now().UnixNano())
        http.HandleFunc("/", handler)
        log.Println("Server running on ", port)
        log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}

~/context$
```

The new dbwriter() function simply sleeps for 0-4 seconds and then writes true to the done channel received in the
parameter list. Our updated handler() function now runs the dbwriter async (as a goroutine) and waits for it to complete
for up to 2 seconds before aborting and reporting a timeout to the user.

Rerun the server:

```
~/context$ go run ctx.go

2022/06/01 04:31:45 Server running on  4444

```

Test the server a few times from another terminal:

```
~$ curl localhost:4444/catamaran

Hit on catamaran!

~$ curl localhost:4444/catamaran

~/context$
```

In the example above, the first request succeeded (the db wrote the request in time) and the second request timed out.
You can try curl with the `-v` switch if you like to see the 504 status returned when the timeout happens:

```
~/context$ curl -v localhost:4444/catamaran

*   Trying 127.0.0.1:4444...
* Connected to localhost (127.0.0.1) port 4444 (#0)
> GET /catamaran HTTP/1.1
> Host: localhost:4444
> User-Agent: curl/7.81.0
> Accept: */*
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 504 Gateway Timeout
< Date: Wed, 01 Jun 2022 04:40:28 GMT
< Content-Length: 0
< 
* Connection #0 to host localhost left intact

~/context$
```

This solution is ok but there is one key weakness, we cannot cancel the db write once started. When the timeout occurs,
we abort back to the client, even though the db write may eventually complete. We could add another `cancel` channel
between the handler and the dbwriter to allow the handler to cancel the dbwriter, however what happens when we need to
add another goroutine from another package that sends a message as part of the user request. Another `cancel` channel?

Further, if we needed to pass network context back to the db writer (like the users identity, auth token, etc.), there's
no easy way to do it without modifying the dbwriter parameter list. Then what of the messaging goroutine we will later
need to add? More parameters? Different ones?

What we really need here is a universal way to pass http headers (call context) to backend goroutines so that the
goroutines can pick and choose the necessary key/value data for themselves. We also need a unified way to cancel one or
many concurrent activities when appropriate. Go context to the rescue!


### 3. Add the ability to cancel request activity on timeout

In this step we'll use Go context to enable the handler to cancel the background db operation when it takes too long.

Update your service as follows:

```go
~/context$ cat ctx.go

package main

import (
        "context"
        "fmt"
        "log"
        "math/rand"
        "net/http"
        "strconv"
        "time"
)

func dbwriter(ctx context.Context, route string, done chan bool) {
        status := false
        defer func() { done <- status }()
        log.Println("db save " + route + " START")
        select {
        case <-time.After(time.Duration(rand.Intn(5)) * time.Second):
                log.Println("db save " + route + " FINISH")
                status = true
        case <-ctx.Done():
                log.Println("db save " + route + " ABORT")
        }
}

func report(status bool, route string, w http.ResponseWriter) {
        if status {
                w.WriteHeader(http.StatusCreated)
                fmt.Fprintf(w, "Hit on %s!\n", route)
        } else {
                w.WriteHeader(http.StatusGatewayTimeout)
                fmt.Fprintf(w, "DB write on %s timed out!\n", route)
        }
}

func handler(w http.ResponseWriter, r *http.Request) {
        route := r.URL.Path[1:]
        done := make(chan bool)
        ctx, cancel := context.WithCancel(context.Background())
        go dbwriter(ctx, route, done)
        select {
        case status := <-done:
                report(status, route, w)
        case <-time.After(2 * time.Second):
                log.Println("Canceling db operation")
                cancel()
                report(<-done, route, w)
        }
}

func main() {
        port := 4444
        rand.Seed(time.Now().UnixNano())
        http.HandleFunc("/", handler)
        log.Println("Server running on ", port)
        log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}

~/context$ 
```

This code updates the dbwriter to take a context.Context object and to use select to monitor the context for
cancelation, via the Done() method. Because the call path is becoming more complex and because the handler() is now
dependent on our write to the done channel, we have moved the done channel write into a defer function. This example
takes a pessimistic approach, assuming things fail (returning status false) unless the execution time case (time.After)
succeeds and sets status to true.

We also update the handler to create the context and pass it to the dbwriter() goroutine, then, if the dbwriter takes
too long, we invoke the cancel() function, telling the dbwriter to shutdown.

Run the new server:

```
~/context$ go run ctx.go

2022/06/01 04:32:23 Server running on  4444

```

Open a new terminal, Now test the server from this client terminal:

```
~$ curl localhost:4444/catamaran

Hit on catamaran!

~$ curl localhost:4444/catamaran

Hit on catamaran!

~$ curl localhost:4444/catamaran

DB write on catamaran timed out!

~$ curl localhost:4444/catamaran

DB write on catamaran timed out!

~$ curl localhost:4444/catamaran

Hit on catamaran!

~$
```

Look at the output on the server:

```
~/gosrc$ go run ctx.go

2022/06/01 04:33:26 Server running on  4444

2022/06/01 04:33:31 db save catamaran START
2022/06/01 04:33:32 db save catamaran FINISH

2022/06/01 04:33:33 db save catamaran START
2022/06/01 04:33:34 db save catamaran FINISH

2022/06/01 04:33:35 db save catamaran START
2022/06/01 04:33:37 Canceling db operation
2022/06/01 04:33:37 db save catamaran ABORT

2022/06/01 04:33:38 db save catamaran START
2022/06/01 04:33:40 Canceling db operation
2022/06/01 04:33:40 db save catamaran ABORT

2022/06/01 04:33:41 db save catamaran START
2022/06/01 04:33:43 db save catamaran FINISH
2022/06/01 04:33:43 Canceling db operation

```

In the example run above you can see some interesting results.

- two of the calls complete normally
- two of the calls timeout normally
- one of the call times out but before the cancel message is received the dbwriter completes its operation

This last case is a little tricky but illustrates a "messages crossing in the mail" case. The handler() requested a
cancel but the dbwriter completed its operation before the cancel arrived. This is an important lesson to take away.
Just because you ask to cancel something, does not mean it was actually canceled. Whenever possible you should wait for
the status of the canceled operation before making assumptions about its state.

Here's the code from our handler that waits to see what really happened after the cancel is sent:

```go
        select {
        case status := <-done:
                report(status, route, w)
        case <-time.After(2 * time.Second):
                log.Println("Canceling db operation")
                cancel()
                report(<-done, route, w)
        }
```

The last request in the above example timed out but the dbwriter completed the operation before it could cancel, writing
`true` to the done channel rather than false. This caused the handler to return 200 OK rather than 504. Note that the
report() function will not be called until a status bool is read from the done channel in the cancel case. This
rendezvous style channel synchronization ensures that the dbwriter go routine has completed or canceled before the
handler reports status back to the user.


### 4. Passing network context through Go context

Imagine our user logged in before accessing our service. When logging in our login host set a GROUP_ID cookie, which
their browser now returns to us with every request in the `cookie` header. Further imagine that the dbwriter can save
the groupid with the route information, enriching the saved data.

Cookies are http headers. JavaScript code running in a browser can not set a cookie. Cookies are part of the underlying
communications layer handling authentication, cache control and many other aspects of the session, all outside the scope
of pure application logic. Propagating this information within our go service is often beneficial but updating function
parameter lists with evolving header metadata is problematic.

What we need is a way to pass arbitrary network context within our Go code. Once again, context.Context to the rescue.
Let's update our program to pass cookies (if present) to all of our backend goroutines via the Context.

Edit your program as follows:

```go
~/context$ vi ctx.go

package main

import (
        "context"
        "fmt"
        "log"
        "math/rand"
        "net/http"
        "strconv"
        "time"
)

type ContextKey string

func dbwriter(ctx context.Context, route string, done chan bool) {
        status := false
        defer func() { done <- status }()
        log.Println("db save " + route + " START")
        select {
        case <-time.After(time.Duration(rand.Intn(5)) * time.Second):
                if group := ctx.Value(ContextKey("GROUP_ID")); group != nil {
                        log.Println("db save", route, group, "FINISH")
                } else {
                        log.Println("db save " + route + " FINISH")
                }
                status = true
        case <-ctx.Done():
                log.Println("db save " + route + " ABORT")
        }
}

func report(status bool, route string, w http.ResponseWriter) {
        if status {
                w.WriteHeader(http.StatusCreated)
                fmt.Fprintf(w, "Hit on %s!\n", route)
        } else {
                w.WriteHeader(http.StatusGatewayTimeout)
                fmt.Fprintf(w, "DB write on %s timed out!\n", route)
        }
}

func handler(w http.ResponseWriter, r *http.Request) {
        route := r.URL.Path[1:]
        done := make(chan bool)
        ctx, cancel := context.WithCancel(context.Background())
        for _, cookie := range r.Cookies() {
                ctx = context.WithValue(ctx, ContextKey(cookie.Name), cookie.Value)
        }
        go dbwriter(ctx, route, done)
        select {
        case status := <-done:
                report(status, route, w)
        case <-time.After(2 * time.Second):
                log.Println("Canceling db operation")
                cancel()
                report(<-done, route, w)
        }
}

func main() {
        port := 4444
        rand.Seed(time.Now().UnixNano())
        http.HandleFunc("/", handler)
        log.Println("Server running on ", port)
        log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}

~/context$
```

The key changes to our code are as follows:

```go
        ctx, cancel := context.WithCancel(context.Background())
        for _, cookie := range r.Cookies() {
                ctx = context.WithValue(ctx, ContextKey(cookie.Name), cookie.Value)
        }
```

In the handler code above, we wrap the original context in a new context with each cookie key/value pair. Note that it
is advisable to use a custom key type when adding values to context. Because context k/v data is essentially stored in a
map under the covers, unique key types ensure that multiple packages can save and retrieve context values independently.

```go
                if group := ctx.Value(ContextKey("GROUP_ID")); group != nil {
                        log.Println("db save", route, group, "FINISH")
                } else {
                        log.Println("db save " + route + " FINISH")
                }
```

The dbwriter code above now attempts to retrieve the `GROUP_ID` from the context. If present we print it out with the
FINISH message.

Let's test it! Run the server:

```
~/context$ go run ctx.go
2022/06/01 04:43:18 Server running on  4444

```

Now run some tests from the client terminal:

```
~$ curl localhost:4444/catamaran

DB write on catamaran timed out!

~$ curl localhost:4444/catamaran

Hit on catamaran!

~$ curl localhost:4444/catamaran --cookie "GROUP_ID=multihulls"

DB write on catamaran timed out!

~$ curl localhost:4444/catamaran --cookie "GROUP_ID=multihulls"

Hit on catamaran!

~$ curl localhost:4444/catamaran --cookie "GROUP_ID=fastboats"

Hit on catamaran!

~$
```

Switch back to the server to make sure the cookies are being delivered to the dbwriter by context:

```
~/gosrc$ go run ctx.go

2022/06/01 04:43:42 Server running on  4444

2022/06/01 04:43:50 db save catamaran START
2022/06/01 04:43:52 Canceling db operation
2022/06/01 04:43:52 db save catamaran ABORT

2022/06/01 04:43:53 db save catamaran START
2022/06/01 04:43:54 db save catamaran FINISH

2022/06/01 04:44:36 db save catamaran START
2022/06/01 04:44:38 Canceling db operation
2022/06/01 04:44:38 db save catamaran ABORT

2022/06/01 04:44:39 db save catamaran START
2022/06/01 04:44:40 db save catamaran multihulls FINISH

2022/06/01 04:47:32 db save catamaran START
2022/06/01 04:47:34 db save catamaran fastboats FINISH
2022/06/01 04:47:34 Canceling db operation

```

Perfect! We have our server passing network call based context into background goroutines.


### 5. Challenge

You may have noticed that our current use case for context invokes the cancel() function after a deadline is reached (2
seconds). The code used above would be a good fit if we wanted to offer an API call allowing the user to cancel this
activity. However, if we are just looking for a timeout, perhaps a deadline style context would be a better fit. Update
your server to use a deadline to cancel background activity after two seconds rather than a cancel() function.


<br>

Congratulations you have completed the lab!!

<br>

