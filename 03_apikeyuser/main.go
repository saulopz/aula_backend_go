package main

import (
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/rs/cors"
)

func main() {
	router := http.NewServeMux()

	router.HandleFunc("GET /auth", func(w http.ResponseWriter, r *http.Request) {
		apiKeys := getAPIKeys("apikeys.json")

		user := r.Header.Get("X-User")
		key := r.Header.Get("X-API-KEY")

		esperada, existe := apiKeys[user]
		if !existe || !chaveValida(key, esperada) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		fmt.Fprintln(w, "Autenticação OK! Olá,", user)
	})

	c := cors.AllowAll()
	handler := c.Handler(router)

	log.Println("Servidor ouvindo em http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func getAPIKeys(arquivo string) map[string]string {
	apiKeys := map[string]string{}

	dados, err := os.ReadFile(arquivo)
	if err != nil {
		return apiKeys
	}

	if err := json.Unmarshal(dados, &apiKeys); err != nil {
		return map[string]string{}
	}
	return apiKeys
}

func chaveValida(recebida, esperada string) bool {
	if recebida == "" || esperada == "" {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(recebida), []byte(esperada)) == 1
}
