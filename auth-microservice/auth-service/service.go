package main

import (
	"database/sql"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func RegisterUser(db *sql.DB, email, password string) error {
	hashed, err := HashPassword(password)
	if err != nil {
		return err
	}
	_, err = db.Exec("INSERT INTO users (email, password) VALUES ($1, $2)", email, hashed)
	return err
}

func LoginUser(db *sql.DB, email, password string) (*User, error) {
	var user User
	err := db.QueryRow("SELECT id, email, password, role, created_at FROM users WHERE email=$1", email).
		Scan(&user.ID, &user.Email, &user.Password, &user.Role, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, err
	}
	return &user, nil
}
