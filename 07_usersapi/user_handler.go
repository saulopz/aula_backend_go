package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

type UserHandler struct {
	repo UserRepository
}

func NewUserHandler(repo UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

func lerUser(w http.ResponseWriter, r *http.Request) (User, bool) {
	var u User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return User{}, false
	}
	if u.Name == "" || u.Email == "" {
		http.Error(w, "Os campos name e email são obrigatórios", http.StatusBadRequest)
		return User{}, false
	}
	return u, true
}

func lerID(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return 0, false
	}
	return id, true
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	u, ok := lerUser(w, r)
	if !ok {
		return
	}

	criado, err := h.repo.Create(u)
	if err != nil {
		erroInterno(w, err)
		return
	}

	w.Header().Set("Location", "/users/"+strconv.Itoa(criado.ID))
	responderJSON(w, http.StatusCreated, criado)
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	users, err := h.repo.List()
	if err != nil {
		erroInterno(w, err)
		return
	}
	responderJSON(w, http.StatusOK, users)
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, ok := lerID(w, r)
	if !ok {
		return
	}

	u, err := h.repo.GetByID(id)
	if errors.Is(err, ErrUserNotFound) {
		http.Error(w, "Usuário não encontrado", http.StatusNotFound)
		return
	}
	if err != nil {
		erroInterno(w, err)
		return
	}
	responderJSON(w, http.StatusOK, u)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := lerID(w, r)
	if !ok {
		return
	}
	u, ok := lerUser(w, r)
	if !ok {
		return
	}
	u.ID = id

	err := h.repo.Update(u)
	if errors.Is(err, ErrUserNotFound) {
		http.Error(w, "Usuário não encontrado", http.StatusNotFound)
		return
	}
	if err != nil {
		erroInterno(w, err)
		return
	}
	responderJSON(w, http.StatusOK, u)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := lerID(w, r)
	if !ok {
		return
	}

	err := h.repo.Delete(id)
	if errors.Is(err, ErrUserNotFound) {
		http.Error(w, "Usuário não encontrado", http.StatusNotFound)
		return
	}
	if err != nil {
		erroInterno(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
