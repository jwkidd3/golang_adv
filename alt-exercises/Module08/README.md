# Using Jwt for user Authentication

### Overview

In this lab you will learn how to use jwt package for generating user tokens and validating them in incoming requests

if you haven't completed the previous exercise you can use the code from the Start folder as your starting point. 

at the end of this module , your project should look like that:

**from jwt.io** : *JSON Web Token (JWT) is an open standard ([RFC 7519](https://tools.ietf.org/html/rfc7519)) that defines a compact and self-contained way for securely transmitting information between parties as a JSON object. This information can be verified and trusted because it is digitally signed. JWTs can be signed using a secret (with the **HMAC** algorithm) or a public/private key pair using **RSA** or **ECDSA**.*

├── cmd

│  └── main.go

├── configs

│  └── app.env

├── go.mod

├── go.sum

├── internal

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

### Start by running mysql 

1. make sure you sql server is still running , if not you can use the following command

```bash
docker run -p 3306:3306 -p 33060:33060 --name mysqldb -v ~/mysql:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=pass -d mysql

```

### Implement UserAuth interface

1. we continuing working with interfaces to be able to support dependency injection easily,

   in the start directory under the internal/users/user.go a new interface was added called UserAuth

```bash
type UserAuth interface {
	IsTokenExists(r *http.Request) (bool, string)
	IsUserTokenValid(token string) bool
	UserFromToken(tokenString string) (*User, error)
	GetTokenForUser(user *User) (string, error)
}
```

2. add a new package **auth** under the users directory and inside it add the **auth.go**  and **token.go** as provided for you in the stater project in this current directory.
3. in the **auth.go** create a type 

```bash
type JwtAuthenticator struct{}
```

4. the JwtAuthenticator type should implement the **UserAuth** interface 

5. export a getter function in **auth.go** the  to get the auth service (this function is provided for you in the auth.go in the starter project)

```bash
var authenticator *JwtAuthenticator

func GetAuthenticator() *JwtAuthenticator {
	if authenticator == nil {
		authenticator = &JwtAuthenticator{}
		initKeys()
	}
	return authenticator
}
```

please not that a skeleton file is provided for you and the public and private keys are provided as well in the starter folder, make sure you add those under your root project directory/keys 

### Create a Jwt middleware to verify incoming requests

1. when a call is being  made to one of the endopints in your app you should verify that the user credentials are allowed to execute the call.

   off course that we do not like to add individual checks for each and every endpoint , one of the best options that we have is using a middleware function to filter out the requests and to allow the calls only to authenticated users.

   call the middleware 

```bash
func (jwtAuth JwtAuthenticator) JwtVerify(next http.Handler) http.Handler
```

2. after a successful token validation the jwt middleware should add to the request context the user object that was validated

   **Tip** you should use the 

```bash
context.WithValue(r.Context(), "user", usr) 
and then call 
next.ServeHTTP(w, r.WithContext(ctx))
```

### Update the UserService

1. update UsersService in user.service.go file  to hold another member of the UserAuth interface , and initialize it with  auth.GetAuthenticator()

```bash
type UsersService struct {
	DB      users.UserDatastore
	JwtAuth users.UserAuth
}
```

2. make sure all the suth calls are now happening only through the users.UserAuth interface
3. change the Login and CreateUser(which actually maps to register) methods to return upon success a serialized json using a generic map, it contains the following:
   1. Status - indication of success
   2. access-token - the actual token 
   3. the details of user that loged in

```bash
var resp = map[string]interface{}{"status": true, "access-token": tokenString, "user": currUser}
	json.NewEncoder(w).Encode(resp)
}
```

### Update the Handler functions routes

1. change the root path for all the user related endpoints (except login and register) , have it start with /auth/

   so it should look like 

   1. /register -> usersService.CreateUser using method POST
   2. /login -> usersService.Login using method POST
   3. /auth/user -> usersService.FetchUsers  using method GET
   4. /auth/user/{id} -> usersService.GetUser  using method GET
   5. /auth/user/{id} -> usersService.UpdateUser using method PUT
   6. /auth/user/{id} -> usersService.DeleteUser using method DELETE

**Tip**  you can use the subrouter feature 

```bash
s := r.PathPrefix("/auth").Subrouter()
```

2. add the Jwt middleware you created to intercept the calls and check the validity of the users 

1. add the Jwt middleware you created to intercept the calls and check the validity of the users 

```bash
s.Use(jv.JwtVerify)
```

### Test 

1. there is another mock that implements the UserAuth interface in the  user.service_test.go set the mocked UserService with it , 
2. make sure test pass



