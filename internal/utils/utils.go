package utils

func Contains[T comparable](targets []T, element T) bool {
	for idx := range targets {
		if targets[idx] == element {
			return true
		}
	}

	return false
}
