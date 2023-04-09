package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/someuser/gameserver/internal/users"
)

// used to hold our user list in memory
var userMap = struct {
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
	fileName := "users.json"
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

	if user, ok := userMap.m[int(userID)]; ok {
		return &user
	}
	return nil
}

func removeUser(userID int) {
	delete(userMap.m, userID)
}

func getUserList() []users.User {
	users := make([]users.User, 0, len(userMap.m))
	for _, value := range userMap.m {
		users = append(users, value)
	}
	return users
}

func getUserIds() []int {
	userIds := []int{}
	for key := range userMap.m {
		userIds = append(userIds, key)
	}
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

	userMap.m[addOrUpdateID] = user

	return addOrUpdateID, nil
}

//CreateUser function -- create a new user
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

//FetchUser function
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
	urlPart := strings.Split(r.URL.Path, fmt.Sprintf("user/"))
	if len(urlPart[1:]) > 1 {
		w.WriteHeader(http.StatusBadRequest)
		return 0
	}
	userID, err := strconv.Atoi(urlPart[len(urlPart)-1])
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusNotFound)
		return 0
	}
	return userID
}
func HandleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		FetchUsers(w, r)
	case http.MethodPost:
		CreateUser(w, r)
		w.WriteHeader(http.StatusCreated)
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func HandleUser(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		GetUser(w, r)
	case http.MethodPut:
		UpdateUser(w, r)
	case http.MethodDelete:
		DeleteUser(w, r)
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
