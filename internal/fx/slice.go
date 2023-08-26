package fx

func MapV[T any, R any](collection []T, fn func(item T) R) []R {
	result := make([]R, len(collection))

	for i, item := range collection {
		result[i] = fn(item)
	}

	return result
}
