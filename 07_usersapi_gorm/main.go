package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// 1. Configuração
	//
	// Lê o arquivo .env da pasta onde o programa foi executado.
	// Se não existir, segue usando as variáveis do sistema.
	if err := godotenv.Load(); err != nil {
		log.Println("Arquivo .env não encontrado, usando as variáveis do sistema")
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("Defina a DATABASE_URL no arquivo .env")
	}

	// 2. Conexão com o banco
	//
	// O gorm.Open já testa a conexão: se o banco estiver fora do ar,
	// ou a senha estiver errada, o erro aparece aqui.
	//
	// logger.Info mostra no terminal cada SQL que o GORM executa.
	// Ótimo para aprender. Em produção, use logger.Warn.
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Não foi possível conectar ao banco: ", err)
	}

	// Por dentro, o GORM mantém um conjunto (pool) de conexões.
	// Pegamos ele só para fechar tudo quando o programa terminar.
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	// 3. Montagem: DAO -> handler
	//
	// Este é o único lugar que sabe que usamos o GORM.
	repo := NewUserGormDAO(db)
	handler := NewUserHandler(repo)

	// 4. Rotas
	router := http.NewServeMux()
	router.HandleFunc("POST /users", handler.Create)
	router.HandleFunc("GET /users", handler.List)
	router.HandleFunc("GET /users/{id}", handler.GetByID)
	router.HandleFunc("PUT /users/{id}", handler.Update)
	router.HandleFunc("DELETE /users/{id}", handler.Delete)

	// CORS liberado para que páginas web possam chamar a API.
	c := cors.AllowAll()

	// 5. Servidor
	log.Println("Servidor ouvindo em http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", c.Handler(router)))
}
