package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func responderJSON(w http.ResponseWriter, status int, dado any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(dado)
}

func erroInterno(w http.ResponseWriter, err error) {
	log.Println("erro interno:", err)
	http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
}
