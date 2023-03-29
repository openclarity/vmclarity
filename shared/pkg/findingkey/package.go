package findingkey

import (
	"fmt"

	"github.com/openclarity/vmclarity/api/models"
)

type PackageKey struct {
	PackageName    string
	PackageVersion string
}

func (k PackageKey) String() string {
	return fmt.Sprintf("%s.%s", k.PackageName, k.PackageVersion)
}

func GeneratePackageKey(info models.PackageFindingInfo) PackageKey {
	return PackageKey{
		PackageName:    *info.Name,
		PackageVersion: *info.Version,
	}
}
