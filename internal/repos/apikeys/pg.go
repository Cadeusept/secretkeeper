package apikeys

import (
	"database/sql"
	"errors"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(dsn string) (*Storage, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return &Storage{db: db}, nil
}

func (s *Storage) AddKey(userID, serviceID string, encryptedKey []byte) error {
	_, err := s.db.Exec(
		"INSERT INTO api_keys (user_id, service_id, api_key) VALUES ($1, $2, $3)",
		userID, serviceID, encryptedKey,
	)
	return err
}

func (s *Storage) GetKey(userID, serviceID string) ([]byte, error) {
	var encryptedKey []byte
	err := s.db.QueryRow(
		"SELECT api_key FROM api_keys WHERE user_id = $1 AND service_id = $2",
		userID, serviceID,
	).Scan(&encryptedKey)
	if err == sql.ErrNoRows {
		return nil, errors.New("key not found")
	}
	return encryptedKey, err
}

func (s *Storage) UpdateKey(userID, serviceID string, encryptedKey []byte) error {
	_, err := s.db.Exec(
		"UPDATE api_keys SET api_key = $1 WHERE user_id = $2 AND service_id = $3",
		encryptedKey, userID, serviceID,
	)
	return err
}
