package routes

import (
	"net/http"
	"net/http/pprof"

	"github.com/gorilla/mux"
	gamesService "github.com/someuser/gameserver/internal/games/service"
	"github.com/someuser/gameserver/internal/users/auth"
	usersService "github.com/someuser/gameserver/internal/users/service"
)

func Handlers() *mux.Router {

	r := mux.NewRouter().StrictSlash(true)

	us := usersService.Get()
	jv := auth.GetAuthenticator()

	r.HandleFunc("/register", us.CreateUser).Methods("POST")
	r.HandleFunc("/login", us.Login).Methods("POST")

	//add support for profiling and tracing of our app
	r.PathPrefix("/debug/pprof/").Handler(http.DefaultServeMux)
	r.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.HandleFunc("/debug/pprof/trace", pprof.Trace)

	s := r.PathPrefix("/auth").Subrouter()
	s.Use(jv.JwtVerify)
	s.HandleFunc("/user", us.FetchUsers).Methods("GET")
	s.HandleFunc("/user/{id}", us.GetUser).Methods("GET")
	s.HandleFunc("/user/{id}", us.UpdateUser).Methods("PUT")
	s.HandleFunc("/user/{id}", us.DeleteUser).Methods("DELETE")
	g := r.PathPrefix("/games").Subrouter()
	g.Use(jv.JwtVerify)
	g.HandleFunc("/gameinfo", gamesService.GetGameInfo).Methods("GET")
	g.HandleFunc("/startnewgame", gamesService.StartNewGame).Methods("GET")
	g.HandleFunc("/joingame/{gametoken}", gamesService.JoinGame).Methods("GET")

	return r
}
