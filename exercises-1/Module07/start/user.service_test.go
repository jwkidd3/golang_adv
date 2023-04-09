package service

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/someuser/gameserver/internal/users"
)

type UserDatastoreMock struct {
	users []users.User
}

func (dbMock *UserDatastoreMock) init() {
	dbMock.users = []users.User{
		{
			ID:       1,
			Name:     "dave",
			Email:    "dave123@gmail.com",
			Password: "dave123",
		},
		{
			ID:       2,
			Name:     "dan",
			Email:    "dan@gmail.com",
			Password: "dan123",
		},
	}
}
func (dbMock *UserDatastoreMock) CreateUser(user *users.User) error {
	//creating the user in the db returns it back with id, it must have all fields in place
	if user.Email != "" && user.Name != "" && user.Password != "" {
		dbMock.users = append(dbMock.users, *user)
		return nil
	}
	return errors.New("missing fields for registring user")
}
func (dbMock *UserDatastoreMock) GetAllUsers() ([]users.User, error) {
	return dbMock.users, nil
}
func (dbMock *UserDatastoreMock) FindUser(email, password string) (*users.User, error) {
	for _, user := range dbMock.users {
		if user.Email == email && user.Password == password {
			return &user, nil
		}
	}
	return nil, errors.New("couldnt find user")
}
func (dbMock *UserDatastoreMock) UpdateUser(id string, user users.User) error {
	for i, u := range dbMock.users {
		if uid, _ := strconv.Atoi(id); uid == int(u.ID) {
			dbMock.users[i] = user
			return nil
		}
	}
	return errors.New("couldnt find user")
}
func (dbMock *UserDatastoreMock) DeleteUser(id string) error {
	for i, u := range dbMock.users {
		if uid, _ := strconv.Atoi(id); uid == int(u.ID) {
			usersSatrt := dbMock.users[:i]
			usersEnd := dbMock.users[i+1:]
			dbMock.users = append(usersSatrt, usersEnd...)
			return nil
		}
	}
	return errors.New("couldnt find user")
}
func (dbMock *UserDatastoreMock) GetUser(id string) (users.User, error) {
	for _, user := range dbMock.users {
		if uid, _ := strconv.Atoi(id); uid == int(user.ID) {
			return user, nil
		}
	}
	return users.User{}, errors.New("couldnt finf user")
}

func TestUsersService_Login(t *testing.T) {

}

func TestUsersService_CreateUser(t *testing.T) {

}

func TestUsersService_FetchUsers(t *testing.T) {

}
