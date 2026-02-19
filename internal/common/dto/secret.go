package dto

const (
	SecretTypeCredentials = "credentials"
	SecretTypeCard        = "card"
	SecretTypeFile        = "file"
	SecretTypeText        = "text"
)

// SecretSupportedTypes доступные типы секретов.
var SecretSupportedTypes = []string{
	SecretTypeCredentials,
	SecretTypeFile,
}

// SecretRequest струкура запроса.
type SecretRequest struct {
	DataType string         `json:"data_type"`
	Name     string         `json:"name"`
	Meta     []MetaData     `json:"meta"`
	EncrData *EncryptedData `json:"data"`
}

// MetaData метаданные секрата.
type MetaData struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// EncryptedData зашифраванные данные секрета.
type EncryptedData struct {
	// Ключ ифрования зашифрованный мастер ключом, закодированн base64
	Key string `json:"key"`
	// Зашифрованные бинарные данные
	Data []byte `json:"data"`
}
