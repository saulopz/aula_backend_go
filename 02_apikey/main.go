package main

import (
	"bufio"
	"crypto/subtle"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/rs/cors"
)

func main() {
	router := http.NewServeMux()

	router.HandleFunc("GET /auth", func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("X-API-KEY")
		if !chaveValida(key, getAPIKey()) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		fmt.Fprintln(w, "Autenticação OK!")
	})

	c := cors.AllowAll()
	handler := c.Handler(router)

	log.Println("Servidor ouvindo em http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func getAPIKey() string {
	file, err := os.Open("apikey.txt")
	if err != nil {
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		return scanner.Text()
	}
	return ""
}

func chaveValida(recebida, esperada string) bool {
	if recebida == "" || esperada == "" {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(recebida), []byte(esperada)) == 1
}