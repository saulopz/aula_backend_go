package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	router := http.NewServeMux()

	router.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello go REST world!")
	})

	log.Println("Servidor ouvindo em http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
