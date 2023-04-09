# Creating the Users Service

### Overview

In this lab you will be creating a users service , this service enables users management using a restful api

if you haven't completed the previous exercise you can use the code from Module01 as your starting point.

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

```
├── service
```

```
│  └── user.service.go
```

```
└── user.go
```

## Getting started

add the users.json file (just a starter file with some users data to get you started with  ) following the directory structure outlined above, (it is included in this foder)

### Start by defining a User Type

1. add a file named user.go following the directory structure outlined above

```bash
package users

//User struct declaration
type User struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
```

### Adding Routes to Support Add,Delete,Get,Update operations

2. change your routes. go file to the following:

```bash
package routes

import (
	"net/http"
	"github.com/<someuser>/gameserver/internal/users/service"
)

func Handlers() {
	usersHandler := http.HandlerFunc(service.HandleUsers)
	userHandler := http.HandlerFunc(service.HandleUser)
	http.Handle("/users", usersHandler)
	http.Handle("/user/", userHandler)
}
```

as you can see there is another package called service under users that holds the buisness logic of the handlers

### Creating the service

1. add a file named users.service.go following the directory structure outlined above , make sure the package name is called service
2. for now we will define an in memory map to hold all the users, to do that you will create a global variable

```bash
// used to hold our user list in memory
var userMap = struct {
	m map[int]users.User
}{m: make(map[int]users.User)}
```

3. implement the handlers based on the requested http method, add the following handlers in users.service.go , and implement the functions that are  called

```bash
func HandleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		FetchUsers(w, r)
	case http.MethodPost:
		CreateUser(w, r)
		w.WriteHeader(http.StatusCreated)
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func HandleUser(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		GetUser(w, r)
	case http.MethodPut:
		UpdateUser(w, r)
	case http.MethodDelete:
		DeleteUser(w, r)
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
```

### Test

make sure that your endpoints return the correct results

you can use curl or postman, to test your endpoints
