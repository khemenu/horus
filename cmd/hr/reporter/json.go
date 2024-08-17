package reporter

import (
	"encoding/json"
	"io"
)

type JsonReporter struct {
	Out io.Writer
}

func (r *JsonReporter) Report(raw any) error {
	return json.NewEncoder(r.Out).Encode(raw)
}
