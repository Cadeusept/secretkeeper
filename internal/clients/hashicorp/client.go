package hashicorp

import (
	"encoding/base64"
	"fmt"
	"os"

	vault "github.com/hashicorp/vault/api"
)

// SecretClient — интерфейс, который ожидает usecase слой
type SecretClient interface {
	Encrypt(plain string) (string, error)
	Decrypt(cipher string) (string, error)
}

// VaultClient реализует SecretClient, используя Vault Transit Engine
type VaultClient struct {
	client *vault.Client
	keyID  string
}

// NewVaultClient создаёт клиента Vault
func NewVaultClient() (*VaultClient, error) {
	// Адрес и токен берутся из переменных окружения
	cfg := vault.DefaultConfig()
	cfg.Address = os.Getenv("VAULT_ADDR") // например: http://localhost:8200

	client, err := vault.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	client.SetToken(os.Getenv("VAULT_TOKEN"))

	keyID := os.Getenv("VAULT_KEY_NAME") // например: "api-key"

	return &VaultClient{
		client: client,
		keyID:  keyID,
	}, nil
}

// Encrypt шифрует строку с помощью Vault Transit
func (vc *VaultClient) Encrypt(plain string) (string, error) {
	base64Plain := base64.StdEncoding.EncodeToString([]byte(plain))

	secret, err := vc.client.Logical().Write(
		fmt.Sprintf("transit/encrypt/%s", vc.keyID),
		map[string]interface{}{
			"plaintext": base64Plain,
		},
	)
	if err != nil {
		return "", err
	}

	cipher, ok := secret.Data["ciphertext"].(string)
	if !ok {
		return "", fmt.Errorf("Vault returned unexpected response")
	}
	return cipher, nil
}

// Decrypt расшифровывает строку с помощью Vault Transit
func (vc *VaultClient) Decrypt(cipher string) (string, error) {
	secret, err := vc.client.Logical().Write(
		fmt.Sprintf("transit/decrypt/%s", vc.keyID),
		map[string]interface{}{
			"ciphertext": cipher,
		},
	)
	if err != nil {
		return "", err
	}

	base64Plain, ok := secret.Data["plaintext"].(string)
	if !ok {
		return "", fmt.Errorf("Vault returned unexpected response")
	}

	decoded, err := base64.StdEncoding.DecodeString(base64Plain)
	if err != nil {
		return "", err
	}
	return string(decoded), nil
}
