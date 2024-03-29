package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	us "github.com/jwkidd3/gameserver/internal/users/service"
)

func Handlers() *mux.Router {

	r := mux.NewRouter().StrictSlash(true)

	r.Use(CommonMiddleware)

	r.HandleFunc("/register", us.CreateUser).Methods("POST")
	r.HandleFunc("/login", us.Login).Methods("POST")
	r.HandleFunc("/user", us.FetchUsers).Methods("GET")
	r.HandleFunc("/user/{id}", us.GetUser).Methods("GET")
	r.HandleFunc("/user/{id}", us.UpdateUser).Methods("PUT")
	r.HandleFunc("/user/{id}", us.DeleteUser).Methods("DELETE")
	return r
}

// CommonMiddleware --Set content-type
func CommonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
		next.ServeHTTP(w, r)
	})
}
