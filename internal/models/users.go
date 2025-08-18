package models

import (
	"database/sql"
	"time"
)

// define new user type
type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

// create new usermodel with wrapped DB connection pool
type UserModel struct {
	DB *sql.DB
}

// use insert method to add new record to users table
func (m *UserModel) Insert(name, email, password string) error {
	return nil
}

// Auth method to verify if user exists with Email/pass, return userID if do
func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	//use exists method to check if user exists with a ID
	return false, nil
}
