package main

import (
	"log"
	"net/http"
	"os"

	"github.com/rs/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("Defina a variável de ambiente DATABASE_URL")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Não foi possível conectar ao banco: ", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	repo := NewUserGormDAO(db)
	// repo := NewUserSQLDAO(sqlDB)
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
