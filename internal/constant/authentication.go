package constant

import (
	"cine-circle/internal/utils"
	"time"
)

const (
	CostHashFunction = 8
	ExpirationDuration = 1 * 24 * time.Hour
	SecretTokenEnv = "SECRET_TOKEN"
	SecretTokenDefault = "secret"
	TokenKind = "Bearer"
	TokenHeader = "Authorization"
	IssToken = "huco-api"
	AuthenticationHeaderName = "Authorization"
	PasswordMinCharacter = 8
	PasswordMaxCharacter = 20
	PasswordAllowUpper = true
	PasswordAllowNumber = true
	PasswordAllowSpecial = true
)

var (
	TokenKey = utils.GetDefaultOrFromEnv(SecretTokenDefault, SecretTokenEnv)
)
