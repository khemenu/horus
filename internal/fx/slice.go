package fx

func FilterV[V any](collection []V, predicate func(item V) bool) []V {
	result := make([]V, 0, len(collection))

	for _, item := range collection {
		if predicate(item) {
			result = append(result, item)
		}
	}

	return result
}

func MapV[T any, R any](collection []T, iteratee func(item T) R) []R {
	result := make([]R, len(collection))

	for i, item := range collection {
		result[i] = iteratee(item)
	}

	return result
}

func FilterMapV[T any, R any](collection []T, callback func(item T) (R, bool)) []R {
	result := make([]R, 0, len(collection))

	for _, item := range collection {
		if r, ok := callback(item); ok {
			result = append(result, r)
		}
	}

	return result
}

func Associate[T any, K comparable, V any](collection []T, transform func(item T) (K, V)) map[K]V {
	result := make(map[K]V, len(collection))

	for _, t := range collection {
		k, v := transform(t)
		result[k] = v
	}

	return result
}
