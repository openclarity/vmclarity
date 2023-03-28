package findingKey

import (
	"fmt"

	"github.com/openclarity/vmclarity/api/models"
)

type RootkitKey struct {
	Name        string
	RootkitType string
	Path        string
}

func (k RootkitKey) String() string {
	return fmt.Sprintf("%s.%s.%s", k.Name, k.RootkitType, k.Path)
}

func GenerateRootkitKey(info models.RootkitFindingInfo) RootkitKey {
	return RootkitKey{
		Name:        *info.RootkitName,
		RootkitType: string(*info.RootkitType),
		Path:        *info.Path,
	}
}
