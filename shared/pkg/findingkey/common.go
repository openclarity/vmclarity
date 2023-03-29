package findingkey

import (
	"fmt"

	"github.com/openclarity/vmclarity/api/models"
)

func GenerateFindingKey(findingInfo *models.Finding_FindingInfo) (string, error) {
	value, err := findingInfo.ValueByDiscriminator()
	if err != nil {
		return "", fmt.Errorf("failed to value by discriminator from finding info: %v", err)
	}

	switch info := value.(type) {
	case models.ExploitFindingInfo:
		return GenerateExploitFindingUniqueKey(info), nil
	case models.VulnerabilityFindingInfo:
		return GenerateVulnerabilityKey(info).String(), nil
	case models.MalwareFindingInfo:
		return GenerateMalwareKey(info).String(), nil
	case models.MisconfigurationFindingInfo:
		return GenerateMisconfigurationKey(info).String(), nil
	case models.RootkitFindingInfo:
		return GenerateRootkitKey(info).String(), nil
	case models.SecretFindingInfo:
		return GenerateSecretKey(info).String(), nil
	case models.PackageFindingInfo:
		return GeneratePackageKey(info).String(), nil
	default:
		return "", fmt.Errorf("unsupported finding info type %T", value)
	}
}
