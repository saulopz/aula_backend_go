package main

import "errors"

var ErrUserNotFound = errors.New("usuário não encontrado")

type UserRepository interface {
	Create(u User) (User, error)
	List() ([]User, error)
	GetByID(id int) (User, error)
	Update(u User) error
	Delete(id int) error
}
