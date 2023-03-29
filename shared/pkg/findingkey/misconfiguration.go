package findingkey

import (
	"fmt"

	"github.com/openclarity/vmclarity/api/models"
)

// MisconfigurationKey One test can report multiple misconfigurations so we need to include the
// message in the unique key.
type MisconfigurationKey struct {
	ScannerName string
	TestID      string
	Message     string
}

func (k MisconfigurationKey) String() string {
	return fmt.Sprintf("%s.%s.%s", k.ScannerName, k.TestID, k.Message)
}

func GenerateMisconfigurationKey(info models.MisconfigurationFindingInfo) MisconfigurationKey {
	return MisconfigurationKey{
		ScannerName: *info.ScannerName,
		TestID:      *info.TestID,
		Message:     *info.Message,
	}
}
