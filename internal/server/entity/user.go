package entity

import "time"

const UserMaxLoginLen = 64

type User struct {
	ID       string    `db:"id"`
	Login    string    `db:"login"`
	Hash     string    `db:"hash"`
	AuthSalt string    `db:"auth_salt"`
	EncrSalt string    `db:"encr_salt"`
	Created  time.Time `db:"created_at"`
}
