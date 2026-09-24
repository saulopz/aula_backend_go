package main

import (
	"encoding/json"
	"net/http"
	"sync"
)

type UserService struct {
	mu   sync.Mutex
	user User
}

func (s *UserService) Get(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	user := s.user
	s.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (s *UserService) Put(w http.ResponseWriter, r *http.Request) {
	var novo User
	if err := json.NewDecoder(r.Body).Decode(&novo); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if novo.Name == "" || novo.Email == "" {
		http.Error(w, "Os campos name e email são obrigatórios", http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	s.user = novo
	s.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(novo)
}
