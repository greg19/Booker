package models

import "database/sql"

const sqlUserTable = `
DROP TABLE IF EXISTS users;
CREATE TABLE users (
 	id       INTEGER PRIMARY KEY AUTOINCREMENT,
 	name 	 TEXT NOT NULL,
 	username TEXT NOT NULL UNIQUE,
 	password TEXT NOT NULL,
 	userType INTEGER NOT NULL
);`

type User struct {
	Id       int
	Name     string
	Username string
	Password string
	UserType int
}

const (
	UserTypeAdmin    = iota // 0
	UserTypeEmployee = iota // 1
	UserTypeCustomer = iota // 2
)

func (u *User) IsAdmin() bool {
	return u.UserType == UserTypeAdmin
}

func (u *User) IsEmployee() bool {
	return u.UserType <= UserTypeEmployee
}

func (u *User) IsCustomer() bool {
	return true
	// or maybe return u.UserType == UserTypeCustomer?
}

func userFromRow(row scannable) (*User, error) {
	var u User
	err := row.Scan(&u.Id, &u.Name, &u.Username, &u.Password, &u.UserType)
	return &u, err
}

const sqlUserByUsername = `
SELECT * FROM users WHERE username = ?`

func GetUserByUsername(db *sql.DB, username string) (*User, error) {
	row := db.QueryRow(sqlUserByUsername, username)
	return userFromRow(row)
}

const sqlUserById = `
SELECT * FROM users WHERE id = ?`

func GetUserById(db *sql.DB, id int) (*User, error) {
	row := db.QueryRow(sqlUserById, id)
	return userFromRow(row)
}

const sqlUserCreate = `
INSERT INTO users (name, username, password, userType) VALUES (?, ?, ?, ?)`

func CreateUser(
	db *sql.DB,
	name string,
	username string,
	password string,
	userType int,
) error {
	_, err := db.Exec(sqlUserCreate, name, username, password, userType)
	return err
}

const sqlUserByType = `
SELECT * FROM users WHERE userType <= ?`

func GetUsersByType(db *sql.DB, userType int) ([]*User, error) {
	rows, err := db.Query(sqlUserByType, userType)
	if err != nil {
		return nil, err
	}
	return readFromRows(rows, userFromRow)
}
