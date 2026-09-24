package main

import (
	"cmp"
	"encoding/json"
	"net/http"
	"slices"
	"strconv"
	"sync"
)

type UserService struct {
	mu     sync.Mutex
	users  map[int]User
	nextID int
}

func NewUserService() *UserService {
	return &UserService{
		users: make(map[int]User),
	}
}

func (s *UserService) Create(w http.ResponseWriter, r *http.Request) {
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
	s.nextID++
	novo.ID = s.nextID
	s.users[novo.ID] = novo
	s.mu.Unlock()

	w.Header().Set("Location", "/users/"+strconv.Itoa(novo.ID))
	responderJSON(w, http.StatusCreated, novo)
}

func (s *UserService) List(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	lista := make([]User, 0, len(s.users))
	for _, user := range s.users {
		lista = append(lista, user)
	}
	s.mu.Unlock()

	slices.SortFunc(lista, func(a, b User) int {
		return cmp.Compare(a.ID, b.ID)
	})

	responderJSON(w, http.StatusOK, lista)
}

func (s *UserService) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	s.mu.Lock()
	user, existe := s.users[id]
	s.mu.Unlock()

	if !existe {
		http.Error(w, "Usuário não encontrado", http.StatusNotFound)
		return
	}

	responderJSON(w, http.StatusOK, user)
}
