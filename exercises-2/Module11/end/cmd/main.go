package main

import (
	"log"
	"net/http"

	"github.com/rs/cors"
	"github.com/someuser/gameserver/internal/routes"
)

func main() {

	r := routes.Handlers()

	http.Handle("/", r)
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"x-access-token"},
	})

	handler := c.Handler(r)

	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		log.Fatal(err)
	}

}
