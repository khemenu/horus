package reporter

import (
	"fmt"
	"io"
)

type TextReporter struct {
	unstructuredReporter
	Out io.Writer
}

func (r *TextReporter) Report(raw any) error {
	_, err := fmt.Fprintln(r.Out, raw)
	return err
}
