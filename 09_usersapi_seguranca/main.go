package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatal("Não foi possível conectar ao banco: ", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	clientes := NewApiClientDAO(db)

	if len(os.Args) == 3 && os.Args[1] == "nova-chave" {
		chave, err := clientes.Create(os.Args[2])
		if err != nil {
			log.Fatal("Não foi possível criar o cliente: ", err)
		}
		fmt.Println("Chave do cliente", os.Args[2]+":")
		fmt.Println(chave)
		fmt.Println("Guarde esta chave agora: ela não será mostrada de novo.")
		return
	}

	repo := NewUserGormDAO(db)
	handler := NewUserHandler(repo)

	api := http.NewServeMux()
	api.HandleFunc("POST /users", handler.Create)
	api.HandleFunc("GET /users", handler.List)
	api.HandleFunc("GET /users/{id}", handler.GetByID)
	api.HandleFunc("PUT /users/{id}", handler.Update)
	api.HandleFunc("DELETE /users/{id}", handler.Delete)

	limitador := NewLimitadorPorCliente(2, 5)

	raiz := http.NewServeMux()
	raiz.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})
	raiz.Handle("/", autenticar(clientes, limitador.Middleware(api)))

	origem := os.Getenv("CORS_ORIGIN")
	if origem == "" {
		origem = "http://localhost:5173"
	}
	c := cors.New(cors.Options{
		AllowedOrigins: []string{origem},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type", "X-API-KEY"},
	})

	servidor := &http.Server{
		Addr:              ":8080",
		Handler:           c.Handler(limitarCorpo(1<<20, raiz)),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	log.Println("Servidor ouvindo em http://localhost:8080")
	log.Fatal(servidor.ListenAndServe())
}
