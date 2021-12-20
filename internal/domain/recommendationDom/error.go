package recommendationDom

import (
	"cine-circle-api/pkg/customError"
)

const (
	errCircleNotFoundCode = "ERR_CIRCLE_NOT_FOUND"
	errMediaNotFoundCode  = "ERR_MEDIA_NOT_FOUND"
)

var (
	errCircleNotFound = customError.NewNotFound().WrapCode(errCircleNotFoundCode)
	errMediaNotFound  = customError.NewNotFound().WrapCode(errMediaNotFoundCode)
)
