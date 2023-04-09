# Dependency injection using interfaces

### Overview

In this lab you will redesign your code seperating data access from the service itself, and learn how use  dependency injection using interfaces in order to be able to inject a mock database for later testing.

you will also learn about the httptest package and they way it help us simplify testing our htto endpoints

if you haven't completed the previous exercise you can use the code from Module06 as your starting point. 

at the end of this module , your project should look like that:

├── cmd

│  └── main.go

├── configs

│  └── app.env

├── go.mod

├── go.sum

└── internal

  ├── routes

  │  └── routes.go

  └── users

​    ├── db

​    │  └── db.go

​    ├── service

​    │  ├── user.service.go

​    │  ├── user.service_test.go

​    │  └── users.service.repo.go

​    └── user.go

## Getting started 

### Start by running mysql 

1. in order get stated quickly we will run mysql database as a docker container,for that use the following command :

```bash
docker run -p 3306:3306 -p 33060:33060 --name mysqldb -v ~/mysql:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=pass -d mysql

```

### Adding datastore interface

1. in your start folder you can see that the user.go file now also contains the following interface 

```bash
type UserDatastore interface {
	CreateUser(user *User) error
	GetAllUsers() ([]User, error)
	FindUser(email, password string) (*User, error)
	UpdateUser(id string, user User) error
	DeleteUser(id string) error
	GetUser(id string) (User, error)
}
```

2. seperate the database operation methods from the user service to another file in the same package called  users.service.repo.go as outlined above

3. add the UsersDB type to the users.service.repo.go and make sure it implements all the UserDatastore method interface above, and 

   ```bash
   type UsersDB struct {
   	*sql.DB
   }
   ```

4. Create a getter function that return the interface only 

```bash

func GetUsersDataStore() users.UserDatastore {
	return &UsersDB{database.Get()}
}
```

### Change the user.service 

1. change the user.service.go file to:
   1. define a type of  UsersService , that holds a users.UserDatastore interface
   2.  export a Get method so that only one instance of the service cn be fetched (singleton)

```bash
type UsersService struct {
	DB users.UserDatastore
}

func Get() *UsersService {

}
```

2. change all the functions in user.service.go to belong to the UserService type
3. make sure that all the database operations are now happening only through the UserDatastore interface

### Changing the Handler function to to use the UserService type

1. for your code to compile ,In  your routes. go file change the handler fnctions to first create a UserService object and then bind it to the routes


### Testing 

in your start directory there is another file called user.service_test.go ,that file needs to be placed in the same folder as the user.service.go.

 in it there are three tetsing functions that you will need you to implement .

```bash
func TestUsersService_Login(t *testing.T) 
func TestUsersService_CreateUser(t *testing.T)
func TestUsersService_FetchUsers(t *testing.T)
```

**Tip** in user.service_test.go there is already a simple mock implementation for the users.UserDatastore that you can use for testing

**don't forget to call init on the mock before using it**

