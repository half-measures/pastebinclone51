package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
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
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	stmt := `INSERT INTO users (name, email, hashed_password, created)
	VALUES(?, ?, ?, UTC_TIMESTAMP())`

	//use Exec() method to insert user details and hashed password into users table
	_, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		//if this returns an error, we use errors.As func to check if
		//its a mysql error. If it does, then error assigned out to mysqlerr var
		//we then check if it equals 1062 and if it does, return dupe email err
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}
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
