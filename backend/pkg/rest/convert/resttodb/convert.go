package resttodb

import (
	"encoding/json"
	"fmt"

	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/backend/pkg/database"
)

func ConvertScanConfig(config *models.ScanConfig) (*database.ScanConfig, error) {
	var ret database.ScanConfig
	var err error

	if config.ScanFamiliesConfig != nil {
		ret.ScanFamiliesConfig, err = json.Marshal(config.ScanFamiliesConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal json: %w", err)
		}
	}

	if config.Scope != nil {
		ret.Scope, err = config.Scope.MarshalJSON()
		if err != nil {
			return nil, fmt.Errorf("failed to marshal json: %w", err)
		}
	}

	if config.Scheduled != nil {
		ret.Scheduled, err = config.Scheduled.MarshalJSON()
		if err != nil {
			return nil, fmt.Errorf("failed to marshal json: %w", err)
		}
	}

	ret.Name = config.Name

	return &ret, nil
}

func ConvertTarget(target *models.Target) (*database.Target, error) {
	disc, err := target.TargetInfo.Discriminator()
	if err != nil {
		return nil, fmt.Errorf("failed to get discriminator: %w", err)
	}
	switch disc {
	case "VMInfo":
		vminfo, err := target.TargetInfo.AsVMInfo()
		if err != nil {
			return nil, fmt.Errorf("failed to convert target to vm info: %w", err)
		}
		return &database.Target{
			Type:             vminfo.ObjectType,
			Location:         *vminfo.Location,
			InstanceID:       *vminfo.InstanceID,
			InstanceProvider: string(*vminfo.InstanceProvider),
		}, nil
	default:
		return nil, fmt.Errorf("unknown target type: %v", disc)
	}
}

// nolint:cyclop
func ConvertScanResult(result *models.TargetScanResult) (*database.ScanResult, error) {
	var ret database.ScanResult
	var err error

	ret.ScanID = result.ScanId
	ret.TargetID = result.TargetId

	if result.Exploits != nil {
		ret.Exploits, err = json.Marshal(result.Exploits)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal json: %w", err)
		}
	}
	if result.Malware != nil {
		ret.Malware, err = json.Marshal(result.Malware)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal json: %w", err)
		}
	}
	if result.Misconfigurations != nil {
		ret.Misconfigurations, err = json.Marshal(result.Misconfigurations)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal json: %w", err)
		}
	}
	if result.Rootkits != nil {
		ret.Rootkits, err = json.Marshal(result.Rootkits)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal json: %w", err)
		}
	}
	if result.Sboms != nil {
		ret.Sboms, err = json.Marshal(result.Sboms)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal json: %w", err)
		}
	}

	if result.Secrets != nil {
		ret.Secrets, err = json.Marshal(result.Secrets)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal json: %w", err)
		}
	}
	if result.Status != nil {
		ret.Status, err = json.Marshal(result.Status)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal json: %w", err)
		}
	}
	if result.Vulnerabilities != nil {
		ret.Vulnerabilities, err = json.Marshal(result.Vulnerabilities)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal json: %w", err)
		}
	}

	return &ret, nil
}

func ConvertScan(scan *models.Scan) (*database.Scan, error) {
	var ret database.Scan
	var err error

	ret.ScanConfigID = scan.ScanConfigId

	ret.ScanEndTime = scan.EndTime

	ret.ScanStartTime = scan.StartTime

	if scan.ScanFamiliesConfig != nil {
		ret.ScanFamiliesConfig, err = json.Marshal(scan.ScanFamiliesConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal json: %w", err)
		}
	}

	if scan.TargetIDs != nil {
		ret.TargetIDs, err = json.Marshal(scan.TargetIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal json: %w", err)
		}
	}

	return &ret, nil
}
