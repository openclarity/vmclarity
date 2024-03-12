package types

import "fmt"

var ErrInvalidScanTemplate = fmt.Errorf("invalid scan template provided")

func (o *ScanTemplate) Validate() error {
	if len(o.Families) == 0 {
		return ErrInvalidScanTemplate
	}
	if len(o.ScanObjectInputs) == 0 {
		return ErrInvalidScanTemplate
	}

	return nil
}
