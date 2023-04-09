package main

import (
	"log"
	"net/http"

	"github.com/someuser/gameserver/internal/routes"
)

func main() {

	// Handle routes
	r := routes.Handlers()

	err := http.ListenAndServe(":5000", r)
	if err != nil {
		log.Fatal(err)
	}

}
