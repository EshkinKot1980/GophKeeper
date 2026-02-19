package entity

import "time"

type Secret struct {
	ID            uint64    `db:"id"`
	UserID        string    `db:"iser_id"`
	DataType      string    `db:"data_type"`
	Name          string    `db:"name"`
	MetaData      string    `db:"meta_data"`
	EncryptedData []byte    `db:"encrypted_data"`
	EncryptedKey  string    `db:"encrypted_key"`
	Created       time.Time `db:"created_at"`
	Updated       time.Time `db:"updated_at"`
}
