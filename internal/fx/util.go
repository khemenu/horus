package fx

func Addr[T any](v T) *T {
	return &v
}

func Default[T comparable](target *T, v T) {
	var zero T
	if *target == zero {
		*target = v
	}
}

func Fallback[T comparable](first T, second T) T {
	var zero T
	if first != zero {
		return first
	}

	return second
}

func Cond[T any](cond bool, t T, f T) T {
	if cond {
		return t
	} else {
		return f
	}
}
