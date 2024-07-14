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

func Cond[C ~bool, T any](cond C, t T, f T) T {
	if cond {
		return t
	} else {
		return f
	}
}

func And(cs ...bool) bool {
	for _, c := range cs {
		if !c {
			return false
		}
	}

	return true
}

func Or(cs ...bool) bool {
	for _, c := range cs {
		if c {
			return true
		}
	}

	return false
}
