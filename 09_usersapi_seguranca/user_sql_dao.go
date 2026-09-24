package main

import (
	"database/sql"
	"errors"
)

type UserSQLDAO struct {
	db *sql.DB
}

func NewUserSQLDAO(db *sql.DB) *UserSQLDAO {
	return &UserSQLDAO{db: db}
}

func (d *UserSQLDAO) Create(u User) (User, error) {
	err := d.db.QueryRow(
		"INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id",
		u.Name, u.Email,
	).Scan(&u.ID)
	return u, err
}

func (d *UserSQLDAO) List() ([]User, error) {
	rows, err := d.db.Query("SELECT id, name, email FROM users ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]User, 0)
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (d *UserSQLDAO) GetByID(id int) (User, error) {
	var u User
	err := d.db.QueryRow(
		"SELECT id, name, email FROM users WHERE id = $1", id,
	).Scan(&u.ID, &u.Name, &u.Email)
	if errors.Is(err, sql.ErrNoRows) {
		return User{}, ErrUserNotFound
	}
	return u, err
}

func (d *UserSQLDAO) Update(u User) error {
	res, err := d.db.Exec(
		"UPDATE users SET name = $1, email = $2 WHERE id = $3",
		u.Name, u.Email, u.ID,
	)
	if err != nil {
		return err
	}
	return verificarAfetadas(res)
}

func (d *UserSQLDAO) Delete(id int) error {
	res, err := d.db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		return err
	}
	return verificarAfetadas(res)
}

func verificarAfetadas(res sql.Result) error {
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrUserNotFound
	}
	return nil
}
