package fx

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
