package fx

import "golang.org/x/exp/constraints"

func Min[T constraints.Ordered](a T, b T) T {
	if a < b {
		return a
	}

	return b
}

func Max[T constraints.Ordered](a T, b T) T {
	if a > b {
		return a
	}

	return b
}

func Clamp[T constraints.Ordered](v T, lo T, hi T) T {
	v = Min(v, hi)
	v = Max(v, lo)
	return v
}
