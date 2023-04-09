package service

import (
	"net/http"

	"github.com/someuser/gameserver/internal/users"
	"github.com/someuser/gameserver/internal/users/db"
)

var usersDb = db.Get()

func Login(w http.ResponseWriter, r *http.Request) {

}

//helper function for finding the user based on email and password
func FindUser(email, password string) (*users.User, error) {

}

//CreateUser function -- create a new user
func CreateUser(w http.ResponseWriter, r *http.Request) {

}

//FetchUsers function
func FetchUsers(w http.ResponseWriter, r *http.Request) {

}

func UpdateUser(w http.ResponseWriter, r *http.Request) {

}

func DeleteUser(w http.ResponseWriter, r *http.Request) {

}

func GetUser(w http.ResponseWriter, r *http.Request) {

}
