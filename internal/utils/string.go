package utils

func SliceContainsStr(slice []string, str string) bool {
	for _, elem := range slice {
		if elem == str {
			return true
		}
	}
	return false
}

func StrSlicesHaveAtLeastOneMatch(slice1, slice2 []string) bool {
	for _, elem1 := range slice1 {
		for _, elem2 := range slice2 {
			if elem1 == elem2 {
				return true
			}
		}
	}
	return false
}
