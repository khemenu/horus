package fx

func Must[T any](obj T, err error) T {
	if err != nil {
		panic(err)
	}
	return obj
}

type ErrCollector[T any] struct {
	v   T
	err error
}

func (c *ErrCollector[T]) To(errs *[]error) T {
	if c.err != nil {
		*errs = append(*errs, c.err)
	}

	return c.v
}

func CollectErr[T any](v T, err error) *ErrCollector[T] {
	return &ErrCollector[T]{v, err}
}
