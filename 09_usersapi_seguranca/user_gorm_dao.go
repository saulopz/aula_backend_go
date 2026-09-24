package main

import (
	"errors"

	"gorm.io/gorm"
)

type UserGormDAO struct {
	db *gorm.DB
}

func NewUserGormDAO(db *gorm.DB) *UserGormDAO {
	return &UserGormDAO{db: db}
}

func (d *UserGormDAO) Create(u User) (User, error) {
	u.ID = 0
	err := d.db.Create(&u).Error
	return u, err
}

func (d *UserGormDAO) List() ([]User, error) {
	users := make([]User, 0)
	err := d.db.Order("id").Find(&users).Error
	return users, err
}

func (d *UserGormDAO) GetByID(id int) (User, error) {
	var u User
	err := d.db.First(&u, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return User{}, ErrUserNotFound
	}
	return u, err
}

func (d *UserGormDAO) Update(u User) error {
	res := d.db.Model(&u).Select("Name", "Email").Updates(u)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}

func (d *UserGormDAO) Delete(id int) error {
	res := d.db.Delete(&User{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrUserNotFound
	}
	return nil
}
