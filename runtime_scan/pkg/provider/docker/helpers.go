package docker

import (
	"fmt"
	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/runtime_scan/pkg/provider"
	"os"
)

func convertTags(tags map[string]string) *[]models.Tag {
	ret := make([]models.Tag, 0, len(tags))
	for key, val := range tags {
		ret = append(ret, models.Tag{
			Key:   key,
			Value: val,
		})
	}
	return &ret
}

func getScanConfigFileName(config *provider.ScanJobConfig) string {
	// TODO (adamtagscherer): check if os.TempDir() doesn't create a new unique directory every time it's being called
	return os.TempDir() + config.AssetScanID + "_scanconfig.yaml"
}

func getAssetId(config *provider.ScanJobConfig) (string, error) {
	objectType, err := config.AssetInfo.ValueByDiscriminator()
	if err != nil {
		return "", fmt.Errorf("failed to get asset object type: %w", err)
	}

	switch value := objectType.(type) {
	case *models.ContainerInfo:
		return *value.Id, nil

	case *models.ContainerImageInfo:
		return *value.Id, nil

	default:
		return "", fmt.Errorf("get asset id not implemented for current object type (%s)", objectType)
	}
}

func getAssetName(config *provider.ScanJobConfig) (string, error) {
	objectType, err := config.AssetInfo.ValueByDiscriminator()
	if err != nil {
		return "", fmt.Errorf("failed to get asset object type: %w", err)
	}

	switch value := objectType.(type) {
	case *models.ContainerInfo:
		return *value.ContainerName, nil

	case *models.ContainerImageInfo:
		return *value.Name, nil

	default:
		return "", fmt.Errorf("get asset id not implemented for current object type (%s)", objectType)
	}
}
