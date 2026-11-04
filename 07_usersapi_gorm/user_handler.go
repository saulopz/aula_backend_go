package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

// UserHandler recebe as requisições HTTP e usa o repositório.
//
// Ele cuida só do HTTP: ler a requisição, validar, escolher
// o código de status e responder. Não há nenhuma linha de SQL aqui.
type UserHandler struct {
	repo UserRepository // o contrato, e não o GORM diretamente
}

// NewUserHandler cria o handler com o repositório que ele vai usar.
func NewUserHandler(repo UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

// ---------------------------------------------------------------
// Funções auxiliares
//
// As duas seguem o mesmo padrão: devolvem o valor e um bool.
// Se o bool for false, a função JÁ respondeu o erro ao cliente,
// e quem chamou só precisa fazer "return".
// ---------------------------------------------------------------

// lerUser lê o JSON do corpo da requisição e confere os campos obrigatórios.
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

// lerID pega o {id} do caminho (ex.: /users/7) e converte para int.
//
// Converter para int também é uma proteção: um número não pode
// carregar um pedaço de SQL.
func lerID(w http.ResponseWriter, r *http.Request) (int, bool) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return 0, false
	}
	return id, true
}

// ---------------------------------------------------------------
// Handlers
//
// Todos seguem o mesmo desenho:
//   1. ler e validar a requisição
//   2. chamar o repositório
//   3. tratar o erro esperado (404), depois o inesperado (500)
//   4. responder
// ---------------------------------------------------------------

// Create trata POST /users. Responde 201 com o usuário criado.
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

	// O cabeçalho Location diz onde o novo recurso pode ser encontrado.
	w.Header().Set("Location", "/users/"+strconv.Itoa(criado.ID))
	responderJSON(w, http.StatusCreated, criado)
}

// List trata GET /users. Responde 200 com a lista.
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	users, err := h.repo.List()
	if err != nil {
		erroInterno(w, err)
		return
	}
	responderJSON(w, http.StatusOK, users)
}

// GetByID trata GET /users/{id}. Responde 200 ou 404.
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// 1. ler e validar
	id, ok := lerID(w, r)
	if !ok {
		return
	}

	// 2. chamar o repositório
	u, err := h.repo.GetByID(id)

	// 3. primeiro o erro esperado, depois qualquer outro
	if errors.Is(err, ErrUserNotFound) {
		http.Error(w, "Usuário não encontrado", http.StatusNotFound)
		return
	}
	if err != nil {
		erroInterno(w, err)
		return
	}

	// 4. responder
	responderJSON(w, http.StatusOK, u)
}

// Update trata PUT /users/{id}. Responde 200 ou 404.
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := lerID(w, r)
	if !ok {
		return
	}
	u, ok := lerUser(w, r)
	if !ok {
		return
	}

	// O ID que vale é o do caminho, e não um "id" que venha no JSON.
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

// Delete trata DELETE /users/{id}. Responde 204 (sem corpo) ou 404.
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
