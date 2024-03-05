# Changing the ServerMux to use gorilla/mux 

### Overview

In this lab you will learn how to use the gorrila/mux library to make your life easier when it comes to routing and you code much cleaner. 

if you haven't completed the previous exercise you can use the code from Module03 as your starting point. 

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

go to the GitHub.com/gorilla/mux url , and take a quick look on all the options that you can use when implementing routing with the gorilla/mux library  

### Start by adding the gorilla/mux package

1. add an import for the "github.com/gorilla/mux" in your routing.go file, it should look like that

```bash
import (
	"net/http"

	"github.com/gorilla/mux"
	usersService "github.com/someuser/gameserver/internal/users/service"
)
```



### Changing the Handler function to use gorilla/mux

1. change your routes. go file so it will include the following routes to the following functions
   1. /register -> usersService.CreateUser using method POST
   2. /login -> usersService.Login using method POST
   3. /user -> usersService.FetchUsers  using method GET
   4. /user/{id} -> usersService.GetUser  using method GET
   5. /user/{id} -> usersService.UpdateUser using method PUT
   6. /user/{id} -> usersService.DeleteUser using method DELETE

2. remove the HandleUsers and HandleUser from routes.go

3. change the parsing logic of getting the id of the user , to use the  mux.Vars() method

### Test 

make sure that your endpoints return the correct results

you can use curl or postman, to test your endpoints

