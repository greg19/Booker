package models

import (
	"database/sql"
	"time"
)

const sqlSessionTable = `
DROP TABLE IF EXISTS sessions;
CREATE TABLE sessions (
	token     TEXT PRIMARY KEY,
	userId    INTEGER NOT NULL,
	expiresAt INTEGER NOT NULL
);`

type Session struct {
	Token     string
	UserId    int
	ExpiresAt time.Time
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

func sessionFromRow(row scannable) (*Session, error) {
	var s Session
	var t int64
	err := row.Scan(&s.Token, &s.UserId, &t)
	s.ExpiresAt = time.Unix(t, 0)
	return &s, err
}

const sqlSessionByToken = `
SELECT * FROM sessions WHERE token = ?`

func GetSessionByToken(db *sql.DB, token string) (*Session, error) {
	row := db.QueryRow(sqlSessionByToken, token)
	return sessionFromRow(row)
}

const sqlSessionCreate = `
INSERT INTO sessions (token, userId, expiresAt) VALUES (?, ?, ?)`

func CreateSession(db *sql.DB, token string, userId int, expiresAt time.Time) error {
	_, err := db.Exec(sqlSessionCreate, token, userId, expiresAt.Unix())
	return err
}

const sqlSessionDelete = `
DELETE FROM sessions WHERE token = ?`

func DeleteSession(db *sql.DB, token string) error {
	_, err := db.Exec(sqlSessionCreate, token)
	return err
}
