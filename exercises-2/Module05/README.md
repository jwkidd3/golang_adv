# Connecting to a database

### Overview

In this lab you will learn how to connect ,and run database operations against mysql database. 

if you haven't completed the previous exercise you can use the code from Module04 as your starting point. 

at the end of this module , your project should look like that:

├── README.md

├── cmd

│  ├── __debug_bin

│  └── main.go

├── configs

│  └── db.env

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

​    └── user.go

## Getting started 

### Start by running mysql 

1. in order get stated quickly we will run mysql database as a docker container,for that use the following command :

```bash
docker run -p 3306:3306 -p 33060:33060 --name mysqldb -v ~/mysql:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=pass -d mysql

```

### Adding db configuration file

1. create a file named db.env and place it under configs dir as outlined above , this file should have the following properties:

```bash
MYSQL_DATABASE = users_db  
MYSQL_PASSWORD = pass
MYSQL_USERNAME = root
MYSQL_SERVICE_HOST = localhost
```

### Connect to the database

1. as a starting point your are provided with a file named db.go under users/db following the directory structure outlined above.

2. the db.go file needs to be completed with the following:

   a.	in db.go implement the  getDBConfig method to return the database configuration information from the app.env file and  from environment variables

   **Tip:** you can user the viper library to help you do that 

```bash
func getDBConfig() (username string, password string,
	databasename string, databaseHost string){...}
```

​	  b.	implement the connectDB method to connect to the database and return the *sql.DB 

```bash
func connectDB() *sql.DB{...}
```

### implement the service functions to Create,Read,Update,Delete users

1. your are given the a skelton for the CRUD methods for user service, your task is to fill those methods for operating on the database and returning results.

```bash
func getDBConfig() (username string, password string,
	databasename string, databaseHost string){...}
```

**Tip** when inserting a user password to the database you want it to be secured, for that you can use the "golang.org/x/crypto/bcrypt" package

```bash
bcrypt.GenerateFromPassword(pasword, bcrypt.DefaultCost)
```

when the user logs in you will need compare agianst the credentials, for that you can use 

```bash
bcrypt.CompareHashAndPassword(not encrypted, encrypted password)
```

### Test 

make sure that your endpoints return the correct results

you can use curl or postman, to test your endpoints

