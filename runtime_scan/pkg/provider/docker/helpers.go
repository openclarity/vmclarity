package docker

import (
	"fmt"
	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/runtime_scan/pkg/provider"
	"os"
	"strings"
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

func getScanName(config *provider.ScanJobConfig) (string, error) {
	id, err := getAssetId(config)
	if err != nil {
		return "", fmt.Errorf("failed to get scan id: %w", err)
	}
	return "scan_" + config.ScanID + "_asset_" + strings.Replace(id, ":", "_", -1), nil
}

func getScanConfigFileName(config *provider.ScanJobConfig) (string, error) {
	scanName, err := getScanName(config)
	if err != nil {
		return "", err
	}

	// TODO (adamtagscherer): check if os.TempDir() doesn't create a new unique directory every time it's being called
	return os.TempDir() + scanName + "_scanconfig.yaml", nil
}

func getAssetId(config *provider.ScanJobConfig) (string, error) {
	objectType, err := config.AssetInfo.Discriminator()
	if err != nil {
		return "", fmt.Errorf("failed to get asset object type: %w", err)
	}

	switch objectType {
	case "ContainerInfo":
		containerInfo, err := config.AssetInfo.AsContainerInfo()
		if err != nil {
			return "", fmt.Errorf("failed to get asset id: %w", err)
		}
		return *containerInfo.Id, nil
	case "ContainerImageInfo":
		containerImageInfo, err := config.AssetInfo.AsContainerImageInfo()
		if err != nil {
			return "", fmt.Errorf("failed to get asset id: %w", err)
		}
		return *containerImageInfo.Id, nil
	default:
		return "", fmt.Errorf("get asset id not implemented for current object type (%s)", objectType)
	}
}

func getAssetName(config *provider.ScanJobConfig) (string, error) {
	objectType, err := config.AssetInfo.Discriminator()
	if err != nil {
		return "", fmt.Errorf("failed to get asset object type: %w", err)
	}

	switch objectType {
	case "ContainerInfo":
		containerInfo, err := config.AssetInfo.AsContainerInfo()
		if err != nil {
			return "", fmt.Errorf("failed to get asset name: %w", err)
		}
		return *containerInfo.ContainerName, nil
	case "ContainerImageInfo":
		containerImageInfo, err := config.AssetInfo.AsContainerImageInfo()
		if err != nil {
			return "", fmt.Errorf("failed to get asset name: %w", err)
		}
		return *containerImageInfo.Name, nil
	default:
		return "", fmt.Errorf("get asset id not implemented for current object type (%s)", objectType)
	}
}
