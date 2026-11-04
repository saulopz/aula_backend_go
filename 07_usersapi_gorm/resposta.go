package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// responderJSON envia uma resposta com o código de status e o dado em JSON.
func responderJSON(w http.ResponseWriter, status int, dado any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(dado)
}

// erroInterno trata as falhas inesperadas (o banco caiu, por exemplo).
//
// O erro completo vai para o log do servidor, onde você pode ler.
// O cliente recebe só uma mensagem genérica: os detalhes do banco
// poderiam ajudar alguém a atacar o sistema.
func erroInterno(w http.ResponseWriter, err error) {
	log.Println("erro interno:", err)
	http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
}
