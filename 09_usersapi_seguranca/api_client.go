package main

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

var ErrClienteNaoEncontrado = errors.New("cliente não encontrado")

type ApiClient struct {
	ID        int
	Name      string
	KeyHash   string
	CreatedAt time.Time
}

type ApiClientDAO struct {
	db *gorm.DB
}

func NewApiClientDAO(db *gorm.DB) *ApiClientDAO {
	return &ApiClientDAO{db: db}
}

func (d *ApiClientDAO) Create(nome string) (string, error) {
	chave, err := gerarChave()
	if err != nil {
		return "", err
	}
	cliente := ApiClient{Name: nome, KeyHash: hashChave(chave)}
	if err := d.db.Create(&cliente).Error; err != nil {
		return "", err
	}
	return chave, nil
}

func (d *ApiClientDAO) FindByKey(chave string) (ApiClient, error) {
	var c ApiClient
	err := d.db.Where("key_hash = ?", hashChave(chave)).First(&c).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ApiClient{}, ErrClienteNaoEncontrado
	}
	return c, err
}
