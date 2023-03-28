package findingKey

import (
	"fmt"

	"github.com/openclarity/vmclarity/api/models"
)

type SecretKey struct {
	Fingerprint string
	StartColumn int
	EndColumn   int
}

func (k SecretKey) String() string {
	return fmt.Sprintf("%s.%s.%s", k.Fingerprint, k.StartColumn, k.EndColumn)
}

func GenerateSecretKey(secret models.SecretFindingInfo) SecretKey {
	return SecretKey{
		Fingerprint: *secret.Fingerprint,
		StartColumn: *secret.StartColumn,
		EndColumn:   *secret.EndColumn,
	}
}
