package reporter

type Reporter interface {
	Report(raw any) error
}

type unstructuredReporter interface {
	unstructured()
}

func IsStructured(r Reporter) bool {
	_, ok := r.(unstructuredReporter)
	return !ok
}
