# Adding support for profiling and tracing

### Overview

In this lab you will add the support for profiling and tracing  your webservice.

you will also use the pprof tool to examine the cpu and memory usage of your app

├── README.md

├── cmd

│  └── main.go

├── configs

│  └── app.env

├── go.mod

├── go.sum

├── internal

│  ├── games

│  │  ├── game.go

│  │  ├── gamemanager.go

│  │  ├── gamemessages.go

│  │  ├── gamesession.go

│  │  ├── gamesutils.go

│  │  ├── player.go

│  │  ├── service

│  │  │  └── game.service.go

│  │  └── utils

│  ├── routes

│  │  └── routes.go

│  └── users

│    ├── auth

│    │  ├── auth.go

│    │  └── token.go

│    ├── db

│    │  └── db.go

│    ├── service

│    │  ├── user.service.go

│    │  ├── user.service.repo.go

│    │  └── user.service_test.go

│    └── user.go

└── keys

  ├── app.rsa

  └── app.rsa.pub

## Getting started 

### Adding profiling to your http app

1. to be able to profile your app you will need to:
   1. Import the "net/http/pprof" package
   2.  add routes to your roues.go file to support the profiling routes

```bash
	//add support for profiling and tracing of our app
	r.PathPrefix("/debug/pprof/").Handler(http.DefaultServeMux)
	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/debug/pprof/trace", pprof.Trace)
```

### CPU Profile

1.  run the following command to interactively watch the cpu profile 

```bash
	go tool pprof http://localhost:8080/debug/pprof/profile
```

2. when entering interactive mod call the **top** command you will get something similar to the following:

```bash
Entering interactive mode (type "help" for commands, "o" for options)
(pprof) top
Showing nodes accounting for 10ms, 100% of 10ms total
      flat  flat%   sum%        cum   cum%
      10ms   100%   100%       10ms   100%  runtime.kevent
         0     0%   100%       10ms   100%  runtime.findrunnable
         0     0%   100%       10ms   100%  runtime.mcall
         0     0%   100%       10ms   100%  runtime.netpoll
         0     0%   100%       10ms   100%  runtime.park_m
         0     0%   100%       10ms   100%  runtime.schedule
```

there isnt much information there to help us understand what's going on...

There is a much better way to look at the high-level performance overview - web command, it generates an SVG graph of hot spots and opens it in a web browser, run the following:

```bash
	go tool pprof --web http://localhost:8080/debug/pprof/profile
```

to get a little more data, lets add some load on the app using apache benchmark tool "ab" and run the prof at the same time

```bash
ab -k -c 8 -n 1000000 "http://127.0.0.1:8080/auth/user/1"
```

3. open your web browser and examine the generated svg file

   

### Heap Profile

1.  run the following command to interactively watch the memory profile 

```bash
	go tool pprof http://127.0.0.1:8080/debug/pprof/heap
```

By default it shows the amount of memory currently in-use,But we are more interested in the number of allocated objects. Call pprof with -alloc_objects option:

```bash
	go tool pprof -alloc_objects http://127.0.0.1:8080/debug/pprof/heap
```

```
  flat  flat%   sum%        cum   cum%
    815282 26.11% 26.11%     815282 26.11%  net/textproto.(*Reader).ReadMIMEHeader
    281682  9.02% 35.13%    1623993 52.00%  net/http.(*conn).readRequest
    263180  8.43% 43.56%     263180  8.43%  net/http.Header.Clone
    229379  7.35% 50.90%     229379  7.35%  net/textproto.MIMEHeader.Set (inline)
    223451  7.16% 58.06%     223451  7.16%  net/textproto.MIMEHeader.Add
    196611  6.30% 64.35%     733798 23.50%  net/http.Error
    180232  5.77% 70.12%     191155  6.12%  context.WithCancel
    142001  4.55% 74.67%    1077179 34.49%  github.com/gorilla/mux.(*Router).ServeHTTP
    108746  3.48% 78.15%     108746  3.48%  regexp.(*bitState).reset
    104865  3.36% 81.51%     104865  3.36%  net.(*conn).Read

```

you can see from here that most of the obehcts are allocated by the http.(*conn).readRequest

### Goroutine profile

1.  Goroutine profile dumps the goroutine call stack and the number of running goroutines ,run the following command to interactively watch the goroutine profile 

```bash
	go tool pprof http://127.0.0.1:8080/debug/pprof/goroutine
```

**Question**: how many goroutine are active ?

### Tracing 

1. to support tracing you will need to add another endpoint to your server

```bash
	r.HandleFunc("/debug/pprof/trace", pprof.Trace)
```

2. To collect a 10 second trace we need to issue a request to the endpoint, 

```bash
	 curl localhost:8080/debug/pprof/trace?seconds=10 > trace.out
```

**Tip** , run on parallel the **ab** tool agin so you get some data to see

1. finally lets run the tool and examine the ui

```bash
	go tool trace trace.out
```

*Security note: beware that exposing pprof handlers to the Internet is not advisable. The recommendation is to expose these endpoints on a different http.Server that is only bound to the loopback interface.* [This blog post](http://mmcloughlin.com/posts/your-pprof-is-showing) *discusses the risks and has code samples on how to properly expose pprof handlers.*

