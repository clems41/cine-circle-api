package utils

import (
	"cine-circle/internal/domain"
	"cine-circle/internal/typedErrors"
	"fmt"
	"strconv"
)

// ContainsID return true if value is in slice
func ContainsID(slice []domain.IDType, value domain.IDType) bool {
	for _, elem := range slice {
		if elem == value {
			return true
		}
	}
	return false
}

//StrToID convert string to domain.IDType
func StrToID(str string) (id domain.IDType, err error) {
	integer, err := strconv.Atoi(str)
	if err != nil {
		err = typedErrors.NewApiBadRequestError(err)
		return
	}
	id = domain.IDType(integer)
	return
}

//IDToStr convert domain.IDType to string
func IDToStr(id domain.IDType) string {
	return fmt.Sprintf("%d", id)
}
