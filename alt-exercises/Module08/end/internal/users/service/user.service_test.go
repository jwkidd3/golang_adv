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

type JwtVerifyMock struct{}

func (jwtMock *JwtVerifyMock) IsTokenExists(r *http.Request) (bool, string) {
	return true, "a-mocked-token"
}

func (jwtMock *JwtVerifyMock) IsUserTokenValid(token string) bool {
	return true
}
func (jwtMock *JwtVerifyMock) UserFromToken(tokenString string) (*users.User, error) {
	return nil, nil
}
func (jwtMock *JwtVerifyMock) GetTokenForUser(user *users.User) (string, error) {
	return "a-mocked-token", nil
}

func TestUsersService_Login(t *testing.T) {
	type fields struct {
		DB users.UserDatastore
	}
	type args struct {
		user users.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    users.User
		wantErr bool
	}{
		{name: "LoginUser", fields: fields{DB: &UserDatastoreMock{}},
			args: args{
				user: users.User{
					Email:    "dave123@gmail.com",
					Password: "dave123",
				}},
			want: users.User{
				Email:    "dave123@gmail.com",
				Password: "dave123",
			}, wantErr: false},

		{name: "LoginUser no email", fields: fields{DB: &UserDatastoreMock{}},
			args: args{
				user: users.User{
					Name:     "dave",
					Password: "dave123",
				}},
			want: users.User{
				Name:     "dave",
				Password: "dave123",
			}, wantErr: true},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := &UsersService{
				DB: tt.fields.DB,
			}
			dbMock := tt.fields.DB.(*UserDatastoreMock)
			dbMock.init()
			jsonuser, _ := json.Marshal(tt.args.user)
			req, _ := http.NewRequest("POST", "/login", strings.NewReader(string(jsonuser)))
			rr := httptest.NewRecorder()

			http.HandlerFunc(us.Login).ServeHTTP(rr, req)

			testMap := make(map[string]interface{})
			err := json.Unmarshal([]byte(rr.Body.String()), &testMap)
			t.Log(testMap)

			if err != nil {
				t.Error("error from unmarshal", err)
			}
			if status := rr.Code; status != http.StatusOK && tt.wantErr == false {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
				t.Fail()
			}

			t.Log(testMap)
		})
	}
}

func TestUsersService_CreateUser(t *testing.T) {

	type fields struct {
		DB users.UserDatastore
	}
	type args struct {
		user users.User
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{name: "CreateUser", fields: fields{DB: &UserDatastoreMock{}},
			args: args{
				user: users.User{
					Name:     "moti",
					Email:    "testuser@gmail.com",
					Password: "test123",
				}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			us := &UsersService{
				DB:      tt.fields.DB,
				JwtAuth: &JwtVerifyMock{},
			}
			jsonuser, _ := json.Marshal(tt.args.user)
			req, _ := http.NewRequest("POST", "/register", strings.NewReader(string(jsonuser)))
			rr := httptest.NewRecorder()

			http.HandlerFunc(us.CreateUser).ServeHTTP(rr, req)

			testMap := make(map[string]interface{})
			err := json.Unmarshal([]byte(rr.Body.String()), &testMap)

			if err != nil {
				t.Error("error from unmarshal", err)
			}

			if status := rr.Code; status != http.StatusCreated {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, http.StatusOK)
				t.Fail()
			}
			t.Log(testMap)
		})
	}
}

func TestUsersService_FetchUsers(t *testing.T) {
	us := &UsersService{
		DB:      &UserDatastoreMock{},
		JwtAuth: &JwtVerifyMock{},
	}

	dbMock := us.DB.(*UserDatastoreMock)
	dbMock.init()

	req, _ := http.NewRequest("GET", "/user", nil)
	rr := httptest.NewRecorder()

	http.HandlerFunc(us.FetchUsers).ServeHTTP(rr, req)

	var users []users.User
	err := json.Unmarshal([]byte(rr.Body.String()), &users)
	t.Log(users)
	if err != nil {
		t.Error("error from unmarshal", err)
	}
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
		t.Fail()
	}
	//now check if we got the user we know of
	if !reflect.DeepEqual(users, dbMock.users) {
		t.Errorf("FetchUsers() = %v, want %v", users, dbMock.users)
	}

}
