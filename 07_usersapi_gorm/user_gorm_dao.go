package main

import (
	"errors"

	"gorm.io/gorm"
)

// UserGormDAO acessa a tabela "users" usando o GORM.
//
// Ele guarda apenas a conexão com o banco, nunca os dados.
// Os dados ficam no banco.
type UserGormDAO struct {
	db *gorm.DB
}

// NewUserGormDAO cria o DAO a partir de uma conexão já aberta.
func NewUserGormDAO(db *gorm.DB) *UserGormDAO {
	return &UserGormDAO{db: db}
}

// Create insere um usuário e devolve ele com o ID gerado pelo banco.
func (d *UserGormDAO) Create(u User) (User, error) {
	// Quem gera o ID é o banco. Se o cliente mandou um "id" no JSON,
	// zeramos aqui: com ID zero, o GORM deixa o banco escolher.
	u.ID = 0

	// INSERT INTO users (name, email) VALUES (...) RETURNING id
	err := d.db.Create(&u).Error
	return u, err
}

// List devolve todos os usuários, em ordem de ID.
func (d *UserGormDAO) List() ([]User, error) {
	// make([]User, 0) garante [] no JSON quando a tabela está vazia
	// (sem isso, o JSON mostraria null).
	users := make([]User, 0)

	// SELECT * FROM users ORDER BY id
	// Sem o ORDER BY, o banco não garante ordem nenhuma.
	err := d.db.Order("id").Find(&users).Error
	return users, err
}

// GetByID busca um usuário pelo ID.
func (d *UserGormDAO) GetByID(id int) (User, error) {
	var u User

	// SELECT * FROM users WHERE id = ? ORDER BY id LIMIT 1
	err := d.db.First(&u, id).Error

	// Quando não encontra, o GORM devolve gorm.ErrRecordNotFound.
	// Traduzimos para o NOSSO erro, que é o que o handler conhece.
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return User{}, ErrUserNotFound
	}
	return u, err
}

// Update substitui o nome e o e-mail de um usuário existente.
func (d *UserGormDAO) Update(u User) error {
	// UPDATE users SET name = ?, email = ? WHERE id = ?
	//
	// Model(&u): o GORM usa o u.ID para montar o WHERE.
	// Select(...): ATENÇÃO, sem ele o GORM ignora campos vazios ("")
	// e não os altera. Com o Select, ele grava os dois campos sempre.
	res := d.db.Model(&u).Select("Name", "Email").Updates(u)
	if res.Error != nil {
		return res.Error
	}

	// Nenhuma linha alterada = não existe usuário com esse ID.
	if res.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}

// Delete remove um usuário pelo ID.
func (d *UserGormDAO) Delete(id int) error {
	// DELETE FROM users WHERE id = ?
	// O &User{} só diz ao GORM de qual tabela excluir.
	res := d.db.Delete(&User{}, id)
	if res.Error != nil {
		return res.Error
	}

	// Nenhuma linha excluída = não existe usuário com esse ID.
	if res.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}
