package test

import (
	"cine-circle/internal/constant"
	"github.com/icrowley/fake"
)

func getFakePassword() string {
	return fake.Password(constant.PasswordMinCharacter, constant.PasswordMaxCharacter, constant.PasswordAllowUpper,
		constant.PasswordAllowNumber, constant.PasswordAllowSpecial)
}
