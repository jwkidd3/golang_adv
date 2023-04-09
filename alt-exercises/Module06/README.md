# Connecting to a database

### Overview

In this **quick** lab you will learn how to connect ,and run database operations against mysql database using context

if you haven't completed the previous exercise you can use the code from Module05 as your starting point. 

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

### running database operations using context

1. each of the methods that you implmented in the previous lab should be altered to use context , the reason for using context is very important especially when it comes to operations that can hang.
2. you should change the CRUD methods to use context with timeout so calls will not block in case of a huge load on the database , give  a 10 seconds time out for each call, for example:

```bash
ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
defer cancel()

row := usersDb.QueryRowContext(ctx, "select id,name,email,password from users where email = ?", email)
```

### Test 

make sure that your endpoints still works return the correct results

you can use curl or postman, to test your endpoints

