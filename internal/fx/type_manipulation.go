package fx

func FromPtrOrF[T any](x *T, callback func() T) T {
	if x == nil {
		return callback()
	}

	return *x
}

func CoalesceOr[T comparable](v0 T, v ...T) (result T) {
	if v0 != result {
		return v0
	}

	for _, e := range v {
		if e != result {
			return e
		}
	}

	return v[len(v)-1]
}
