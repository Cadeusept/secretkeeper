package secretkeeper

import (
	"context"
	"log"

	"github.com/Cadeusept/secretkeeper/internal/clients/hashicorp"
	"github.com/Cadeusept/secretkeeper/internal/repos/apikeys"
)

type Usc struct {
	log     log.Logger
	client  hashicorp.SecretClient
	storage apikeys.Storage
}

// AddKey handles the business logic of adding a key
func (u *Usc) AddKey(ctx context.Context, userID, serviceID, apiKey string) error {
	encryptedKey, err := u.client.Encrypt(apiKey)
	if err != nil {
		return err
	}
	return u.storage.AddKey(userID, serviceID, []byte(encryptedKey))
}

// GetKey handles the business logic of retrieving a key
func (u *Usc) GetKey(ctx context.Context, userID, serviceID string) (string, error) {
	encryptedKey, err := u.storage.GetKey(userID, serviceID)
	if err != nil {
		return "", err
	}
	return u.client.Decrypt(string(encryptedKey))
}

// UpdateKey handles updating an existing key
func (u *Usc) UpdateKey(ctx context.Context, userID, serviceID, newKey string) error {
	encryptedKey, err := u.client.Encrypt(newKey)
	if err != nil {
		return err
	}
	return u.storage.UpdateKey(userID, serviceID, []byte(encryptedKey))
}
