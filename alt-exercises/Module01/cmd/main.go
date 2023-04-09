package main

import (
	"log"
	"net/http"

	"github.com/jwkidd3/gameserver/internal/routes"
)

func main() {

	// Handle routes
	routes.Handlers()

	err := http.ListenAndServe(":5001", nil)
	if err != nil {
		log.Fatal(err)
	}

}
