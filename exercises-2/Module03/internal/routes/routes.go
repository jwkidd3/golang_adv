package routes

import (
	"net/http"

	"github.com/someuser/gameserver/internal/users/service"
)

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

func Handlers() {
	usersHandler := http.HandlerFunc(service.HandleUsers)
	userHandler := http.HandlerFunc(service.HandleUser)
	http.Handle("/users", CommonMiddleware(usersHandler))
	http.Handle("/user/", CommonMiddleware(userHandler))
}
