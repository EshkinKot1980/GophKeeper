package dto

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type AuthResponse struct {
	// токен для авторизации (JWT)
	Token string `json:"token"`
	// Соль для создания мастер из пароля ключа закодированная base64
	EncrSalt string `json:"encr_salt"`
}
