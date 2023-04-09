package service

import (
	"database/sql"
	"errors"
	"log"

	"github.com/someuser/gameserver/internal/users"
	database "github.com/someuser/gameserver/internal/users/db"

	"golang.org/x/crypto/bcrypt"
)

type UsersDB struct {
	*sql.DB
}

func GetUsersDataStore() users.UserDatastore {
	return &UsersDB{database.Get()}
}

//CreateUser function -- create a new user
func (db *UsersDB) CreateUser(user *users.User) error {

	if user.Email == "" || user.Password == "" || user.Name == "" {
		return errors.New("cant have empty fields")
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
		return errors.New("Password Encryption failed")
	}
	user.Password = string(pass)

	result, err := db.Exec("insert into users(name,email,password)values(?,?,?)", user.Name, user.Email, user.Password)

	if err != nil {
		return err
	}

	id, e := result.LastInsertId()
	if e != nil {
		return e
	}

	user.ID = uint(id)

	return nil
}

func (db *UsersDB) GetAllUsers() ([]users.User, error) {
	var theUsers []users.User

	rows, err := db.Query("select id,name , email , password  from users")
	defer rows.Close()
	if err != nil {
		log.Print("error occued during users fetch ", err.Error())
		return nil, err
	}

	for rows.Next() {
		var user users.User
		rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password)
		theUsers = append(theUsers, user)
	}
	return theUsers, nil
}

func (db *UsersDB) FindUser(email, password string) (*users.User, error) {
	user := &users.User{}

	if email == "" || password == "" {
		return nil, errors.New("cant have empty email or password")
	}
	row := db.QueryRow("select id,name,email,password from users where email = ?", email)

	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password)

	if err == sql.ErrNoRows {
		return nil, err
	}

	errf := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if errf != nil { //Password does not match!
		return nil, errors.New("Invalid login credentials. Please try again")
	}

	return user, nil
}

func (db *UsersDB) UpdateUser(id string, user users.User) error {

	result, err := db.Exec("update users set name = ? , email= ? ,password = ? where id = ?", user.Name, user.Email, user.Password, id)
	if err != nil {
		log.Print("error occued during user update ", err.Error())
		return err
	}
	num, err := result.RowsAffected()
	if err != nil {
		log.Fatal("couldnt update database ", err.Error())
		return err
	}

	log.Println("number of rows affected is ", num)
	return nil
}

func (db *UsersDB) DeleteUser(id string) error {

	result, err := db.Exec("delete from users where id = ?", id)
	if err != nil {
		log.Print("error occued during user update ", err.Error())
		return err
	}
	_, err = result.RowsAffected()
	if err != nil {
		log.Fatal("couldnt update database ", err.Error())
		return err
	}
	return nil
}

func (db *UsersDB) GetUser(id string) (users.User, error) {

	var user users.User

	row := db.QueryRow("select id,name,email,password from users where id=?", id)
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Password)

	if err == sql.ErrNoRows {
		return users.User{}, err
	}
	return user, nil
}
