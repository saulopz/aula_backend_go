package main

// User é o nosso objeto de transferência: os dados de um usuário.
//
// As tags `json:"..."` dizem como cada campo aparece no JSON.
// Repare que não há nada sobre o banco aqui: o GORM descobre
// a tabela e as colunas pelas convenções de nomes
// (User vira a tabela "users", Name vira a coluna "name"...).
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
