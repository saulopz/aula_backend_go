package main

import (
	"log"
	"net/http"

	"github.com/rs/cors"
)

func main() {
	service := NewUserService()

	router := http.NewServeMux()
	router.HandleFunc("POST /users", service.Create)
	router.HandleFunc("GET /users", service.List)
	router.HandleFunc("GET /users/{id}", service.GetByID)

	c := cors.AllowAll()
	handler := c.Handler(router)

	log.Println("Servidor ouvindo em http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
