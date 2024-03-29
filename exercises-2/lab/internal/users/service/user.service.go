package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
	"github.com/jwkidd3/gameserver/internal/users"
)

// used to hold our user list in memory
var userMap = struct {
	sync.RWMutex
	m map[int]users.User
}{m: make(map[int]users.User)}

func init() {
	fmt.Println("loading users...")
	m, err := loadUserMap()
	userMap.m = m
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d users loaded...\n", len(userMap.m))
}

func loadUserMap() (map[int]users.User, error) {
	fileName := "./cmd/users.json"
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("file [%s] does not exist", fileName)
	}

	file, _ := ioutil.ReadFile(fileName)
	userList := make([]users.User, 0)
	err = json.Unmarshal([]byte(file), &userList)
	if err != nil {
		log.Fatal(err)
	}
	uMap := make(map[int]users.User)
	for i := 0; i < len(userList); i++ {
		uMap[int(userList[i].ID)] = userList[i]
	}
	return uMap, nil
}

func getUser(userID uint) *users.User {
	userMap.RLock()
	defer userMap.RUnlock()
	if user, ok := userMap.m[int(userID)]; ok {
		return &user
	}
	return nil
}

func removeUser(userID int) {
	userMap.Lock()
	defer userMap.Unlock()
	delete(userMap.m, userID)
}

func getUserList() []users.User {
	userMap.RLock()
	users := make([]users.User, 0, len(userMap.m))
	for _, value := range userMap.m {
		users = append(users, value)
	}
	userMap.RUnlock()
	return users
}

func getUserIds() []int {
	userMap.RLock()
	userIds := []int{}
	for key := range userMap.m {
		userIds = append(userIds, key)
	}
	userMap.RUnlock()
	sort.Ints(userIds)
	return userIds
}

func getNextUserID() int {
	userIds := getUserIds()
	return userIds[len(userIds)-1] + 1
}

func addOrUpdateUser(user users.User) (int, error) {
	// if the user id is set, update, otherwise add
	addOrUpdateID := -1
	if user.ID > 0 {
		oldUser := getUser(user.ID)
		// if it exists, replace it, otherwise return error
		if oldUser == nil {
			return 0, fmt.Errorf("user id [%d] doesn't exist", user.ID)
		}
		addOrUpdateID = int(user.ID)
	} else {
		addOrUpdateID = getNextUserID()
		user.ID = uint(addOrUpdateID)
	}
	userMap.Lock()
	userMap.m[addOrUpdateID] = user
	userMap.Unlock()
	return addOrUpdateID, nil
}

func FindUser(email string, password string) (*users.User, error) {
	for _, user := range getUserList() {
		if user.Email == email && password == user.Password {
			return &user, nil
		}
	}
	return nil, errors.New("Invalid login credentials. Please try again")
}

func Login(w http.ResponseWriter, r *http.Request) {
	user := &users.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		var resp = map[string]interface{}{"status": false, "message": "Invalid request"}
		json.NewEncoder(w).Encode(resp)
		return
	}
	currUser, err := FindUser(user.Email, user.Password)

	if err != nil {
		var resp = map[string]interface{}{"status": false, "message": err.Error()}
		json.NewEncoder(w).Encode(resp)
		return
	}

	//update those before sending back
	user.Name = currUser.Name
	user.ID = currUser.ID

	var resp = map[string]interface{}{"status": true, "user": user}
	json.NewEncoder(w).Encode(resp)
}

// CreateUser function -- create a new user
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user users.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = addOrUpdateUser(user)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

// FetchUser function
func FetchUsers(w http.ResponseWriter, r *http.Request) {
	usersList := getUserList()
	j, err := json.Marshal(usersList)
	if err != nil {
		log.Fatal(err)
	}
	_, err = w.Write(j)
	if err != nil {
		log.Fatal(err)
	}
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID := ParseURL(w, r)
	if userID == 0 {
		return
	}
	var user users.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if user.ID != uint(userID) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = addOrUpdateUser(user)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userID := ParseURL(w, r)
	if userID == 0 {
		return
	}
	removeUser(userID)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	userID := ParseURL(w, r)
	if userID == 0 {
		return
	}
	user := getUser(uint(userID))
	if user == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	j, err := json.Marshal(user)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, err = w.Write(j)
	if err != nil {
		log.Fatal(err)
	}
}

func ParseURL(w http.ResponseWriter, r *http.Request) int {

	vars := mux.Vars(r)
	idstr := vars["id"]
	userID, err := strconv.Atoi(idstr)

	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusNotFound)
		return 0
	}
	return userID
}
