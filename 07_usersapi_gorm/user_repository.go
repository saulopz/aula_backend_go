package main

import "errors"

// ErrUserNotFound é o erro que o DAO devolve quando o usuário não existe.
//
// Os handlers comparam com ele (errors.Is) para responder 404.
// Assim, os handlers nunca precisam conhecer os erros do GORM.
var ErrUserNotFound = errors.New("usuário não encontrado")

// UserRepository é o contrato do nosso DAO: as operações de um CRUD.
//
// O handler depende deste contrato, e não do GORM. Qualquer tipo
// que tenha estes cinco métodos pode ser usado no lugar do UserGormDAO.
// (Na página extra "Persistência com SQL" fazemos exatamente isso.)
type UserRepository interface {
	Create(u User) (User, error)
	List() ([]User, error)
	GetByID(id int) (User, error)
	Update(u User) error
	Delete(id int) error
}
