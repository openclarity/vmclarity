package types

import "fmt"

var ErrInvalidScanTemplate = fmt.Errorf("invalid scan template provided")

func (o *ScanTemplate) Validate() error {
	if len(o.Inputs) == 0 {
		return fmt.Errorf("inputs not specified: %w", ErrInvalidScanTemplate)
	}

	return nil
}
