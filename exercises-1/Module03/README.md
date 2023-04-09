# Adding a Middleware

### Overview

In this lab you will just add cors Middleware support for the users service you creted in the previous lab 

if you haven't completed the previous exercise you can use the code from Module02 as your starting point. 

at the end of this module , your project should look like that:

├── README.md

├── cmd

│  ├── main.go

│  └── users.json

├── go.mod

└── internal

  ├── routes

  │  └── routes.go

  └── users

​    ├── service

​    │  └── user.service.go

​    └── user.go

## Getting started

when the browser access a url it adds the "Origin" Header to the request identifieng the domain the request came from. if the upstream service does not Allow this Origin specifically by return an allow header the browser will fail the request.

start by testing the current status of your service when trying to access it from a browser, for that you have in the test folder a javascript inside index.html that you can use to test.

that script calls your service /users endpoint from the browser using ajax .

After opening the index.html file in your browser, right click -> inspect -> console, and see that the request for listing the users is currently failing

### Create the middleware function

1. add the following function to the routes.go file , for demostratoin it will allow any origin *

```bash
// CommonMiddleware --Set content-type
func CommonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
		next.ServeHTTP(w, r)
	})
}
```

2. change your routes accordingly to support middleware by wrapping the handler calls

### Test 

run the index.html file in your browser, right click -> inspect -> console, and see that the request is passing

