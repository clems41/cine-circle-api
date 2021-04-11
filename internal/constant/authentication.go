package constant

import "time"

const (
	Cost = 8
	ExpirationDuration = 1 * 24 * time.Hour
	SecretTokenEnv = "SECRET_TOKEN"
	SecretTokenDefault = "secret"
	TokenKind = "Bearer"
	TokenHeader = "Authorization"
)
