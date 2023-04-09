package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/someuser/gameserver/internal/users"
	"github.com/someuser/gameserver/internal/users/db"
	"golang.org/x/crypto/bcrypt"
)

var usersDb = db.Get()

func Login(w http.ResponseWriter, r *http.Request) {

	user := &users.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	currUser, err := FindUser(user.Email, user.Password)

	if err != nil {
		log.Print("error occued FindUser ", err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var resp = map[string]interface{}{"status": true, "user": currUser}
	json.NewEncoder(w).Encode(resp)
}

func FindUser(email, password string) (*users.User, error) {
	user := &users.User{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	row := usersDb.QueryRowContext(ctx, "select id,name,email,password from users where email = ?", email)

	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password)

	if err == sql.ErrNoRows {
		return nil, errors.New("Email address not found")
	}

	errf := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if errf != nil { //Password does not match!
		return nil, errors.New("Invalid login credentials. Please try again")
	}

	return user, nil
}

//CreateUser function -- create a new user
func CreateUser(w http.ResponseWriter, r *http.Request) {

	user := &users.User{}
	json.NewDecoder(r.Body).Decode(user)

	_, err := FindUser(user.Email, user.Password)

	if err == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user.Password = string(pass)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	result, err := usersDb.ExecContext(ctx, "insert into users(name,email,password)values(?,?,?)", user.Name, user.Email, user.Password)

	if err != nil {
		log.Print("error occued CreateUser ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if id, e := result.LastInsertId(); e != nil {
		log.Println("no rows affected")
		return
	} else {
		user.ID = uint(id)
	}

	w.WriteHeader(http.StatusCreated)

	var resp = map[string]interface{}{"status": true, "user": user}
	json.NewEncoder(w).Encode(resp)
}

//FetchUser function
func FetchUsers(w http.ResponseWriter, r *http.Request) {
	var theUsers []users.User

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	rows, err := usersDb.QueryContext(ctx, "select id,name , email , password  from users")
	defer rows.Close()
	if err != nil {
		log.Print("error occued during FetchUsers ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for rows.Next() {
		var user users.User
		rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
		theUsers = append(theUsers, user)
	}
	json.NewEncoder(w).Encode(theUsers)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	user := &users.User{}
	params := mux.Vars(r)
	var id = params["id"]

	json.NewDecoder(r.Body).Decode(user)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	result, err := usersDb.ExecContext(ctx, "update users set name = ? , email= ? ,password = ? where id = ?", user.Name, user.Email, user.Password, id)
	if err != nil {
		log.Print("error occued during user update ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
	num, err := result.RowsAffected()
	if err != nil {
		log.Fatalf("couldnt update database ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println("number of rows affected is ", num)
	json.NewEncoder(w).Encode(&user)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var id = params["id"]

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	result, err := usersDb.ExecContext(ctx, "delete from users where id = ?", id)
	if err != nil {
		log.Print("error occued during user delete ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = result.RowsAffected()
	if err != nil {
		log.Print("couldnt update database ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode("User deleted")
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var id = params["id"]
	var user users.User

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	row := usersDb.QueryRowContext(ctx, "select id,name,email,password from users where id=?", id)
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password)

	if err != nil {
		log.Print("error occued during user delete ", err.Error())
		w.WriteHeader(http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(&user)
}
