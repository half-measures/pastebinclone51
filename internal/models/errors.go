package models

import "errors"

var (
	ErrNoRecord = errors.New("models: no matching record found")
	// invalid cred err if invalid username or pass
	ErrInvalidCredentials = errors.New("models: invalid credentials")

	// email must be uniq, its a constraint on the mysql table column
	ErrDuplicateEmail = errors.New("models: duplicate email")
)
