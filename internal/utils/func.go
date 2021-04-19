package utils

import "cine-circle/internal/domain"

// ContainsID return true if value is in slice
func ContainsID(slice []domain.IDType, value domain.IDType) bool {
	for _, elem := range slice {
		if elem == value {
			return true
		}
	}
	return false
}
