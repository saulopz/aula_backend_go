package main

import (
	"log"
	"net/http"

	"github.com/rs/cors"
)

func main() {
	service := &UserService{}

	router := http.NewServeMux()
	router.HandleFunc("GET /user", service.Get)
	router.HandleFunc("PUT /user", service.Put)

	c := cors.AllowAll()
	handler := c.Handler(router)

	log.Println("Servidor ouvindo em http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
