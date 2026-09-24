package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/rs/cors"
)

func main() {
	var mu sync.Mutex
	number := 0

	router := http.NewServeMux()

	router.HandleFunc("GET /number", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		valor := number
		mu.Unlock()

		fmt.Fprintln(w, valor)
	})

	router.HandleFunc("PUT /number", func(w http.ResponseWriter, r *http.Request) {
		corpo, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Erro ao ler a requisição", http.StatusBadRequest)
			return
		}

		novo, err := strconv.Atoi(strings.TrimSpace(string(corpo)))
		if err != nil {
			http.Error(w, "O valor precisa ser um número inteiro", http.StatusBadRequest)
			return
		}

		mu.Lock()
		number = novo
		mu.Unlock()

		fmt.Fprintln(w, "number atualizado para", novo)
	})

	c := cors.AllowAll()
	handler := c.Handler(router)

	log.Println("Servidor ouvindo em http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
