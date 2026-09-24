package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/cors"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("Defina a variável de ambiente DATABASE_URL")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Não foi possível conectar ao banco: ", err)
	}

	repo := NewUserSQLDAO(db)
	handler := NewUserHandler(repo)

	router := http.NewServeMux()
	router.HandleFunc("POST /users", handler.Create)
	router.HandleFunc("GET /users", handler.List)
	router.HandleFunc("GET /users/{id}", handler.GetByID)
	router.HandleFunc("PUT /users/{id}", handler.Update)
	router.HandleFunc("DELETE /users/{id}", handler.Delete)

	c := cors.AllowAll()

	log.Println("Servidor ouvindo em http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", c.Handler(router)))
}
