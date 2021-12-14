package exempleDom

import "cine-circle-api/pkg/customError"

const (
	errExempleIntrouvableCode = "ERR_EXEMPLE_INTROUVABLE"
)

var (
	errExempleIntrouvable = customError.NewNotFound().WrapCode(errExempleIntrouvableCode)
)
