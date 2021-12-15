package mediaDom

import "cine-circle-api/pkg/customError"

const (
	errMediaNotFoundCode = "ERR_MEDIA_NOT_FOUND"
)

var (
	errMediaNotFound = customError.NewNotFound().WrapCode(errMediaNotFoundCode)
)
