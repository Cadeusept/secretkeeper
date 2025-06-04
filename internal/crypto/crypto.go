package crypto

import (
	"log"

	"github.com/hashicorp/vault/api"
	"github.com/pedroalbanese/kuznechik"
)

var secretKey []byte

// Инициализация клиента HashiCorp Vault
func initVault(addr, token, secretPath, keyName string) error {
	client, err := api.NewClient(&api.Config{Address: addr})
	if err != nil {
		return err
	}

	client.SetToken(token)

	secret, err := client.Logical().Read(secretPath)
	if err != nil {
		return err
	}

	if secret == nil || secret.Data[keyName] == nil {
		return err
	}

	secretKey = []byte(secret.Data[keyName].(string))
	if len(secretKey) != 32 {
		return err
	}

	log.Println("Vault encryption key loaded successfully")
	return nil
}

func Encrypt(data string) ([]byte, error) {
	block, err := kuznechik.NewCipher(secretKey)
	if err != nil {
		return nil, err
	}
	encrypted := make([]byte, len(data))
	block.Encrypt(encrypted, []byte(data))
	return encrypted, nil
}

func Decrypt(data []byte) (string, error) {
	block, err := kuznechik.NewCipher(secretKey)
	if err != nil {
		return "", err
	}
	decrypted := make([]byte, len(data))
	block.Decrypt(decrypted, data)
	return string(decrypted), nil
}
