package routes

import (
	"net/http"

	"github.com/someuser/gameserver/internal/users/service"
)

func Handlers() {
	usersHandler := http.HandlerFunc(service.HandleUsers)
	userHandler := http.HandlerFunc(service.HandleUser)
	http.Handle("/users", usersHandler)
	http.Handle("/user/", userHandler)
}
