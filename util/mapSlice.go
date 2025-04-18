package util

func MapSlice[T, U any](slice []T, callback func(T) U) []U {
	result := make([]U, cap(slice))

	for index, elem := range slice {
		result[index] = callback(elem)
	}

	return result
}
