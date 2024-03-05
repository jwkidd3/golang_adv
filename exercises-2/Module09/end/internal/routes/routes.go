package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	gamesService "github.com/someuser/gameserver/internal/games/service"
	"github.com/someuser/gameserver/internal/users/auth"
	usersService "github.com/someuser/gameserver/internal/users/service"
)

func Handlers() *mux.Router {

	r := mux.NewRouter().StrictSlash(true)

	r.Use(CommonMiddleware)

	us := usersService.Get()
	jv := auth.GetAuthenticator()

	r.HandleFunc("/register", us.CreateUser).Methods("POST")
	r.HandleFunc("/login", us.Login).Methods("POST")

	s := r.PathPrefix("/auth").Subrouter()
	s.Use(jv.JwtVerify)
	s.HandleFunc("/user", us.FetchUsers).Methods("GET")
	s.HandleFunc("/user/{id}", us.GetUser).Methods("GET")
	s.HandleFunc("/user/{id}", us.UpdateUser).Methods("PUT")
	s.HandleFunc("/user/{id}", us.DeleteUser).Methods("DELETE")
	g := r.PathPrefix("/games").Subrouter()
	g.Use(CommonMiddleware)
	g.Use(jv.JwtVerify)
	g.HandleFunc("/gameinfo", gamesService.GetGameInfo).Methods("GET")
	g.HandleFunc("/startnewgame", gamesService.StartNewGame).Methods("GET")
	g.HandleFunc("/joingame/{gametoken}", gamesService.JoinGame).Methods("GET")

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
