package db_to_rest

import (
	"encoding/json"
	"fmt"

	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/backend/pkg/database"
	"github.com/openclarity/vmclarity/runtime_scan/pkg/utils"
)

func ConvertScanConfig(config *database.ScanConfig) (*models.ScanConfig, error) {
	var ret models.ScanConfig

	if config.ScanFamiliesConfig != nil {
		ret.ScanFamiliesConfig = &models.ScanFamiliesConfig{}
		if err := json.Unmarshal(config.ScanFamiliesConfig, ret.ScanFamiliesConfig); err != nil {
			return nil, err
		}
	}

	if config.Scope != nil {
		ret.Scope = &models.ScanScopeType{}
		if err := ret.Scope.UnmarshalJSON(config.Scope); err != nil {
			return nil, err
		}
	}

	if config.Scheduled != nil {
		ret.Scheduled = &models.RuntimeScheduleScanConfigType{}
		if err := ret.Scheduled.UnmarshalJSON(config.Scheduled); err != nil {
			return nil, err
		}
	}

	ret.Name = config.Name

	ret.Id = utils.StringPtr(config.ID.String())

	return &ret, nil
}

func ConvertScanConfigs(configs []*database.ScanConfig, total int64) (*models.ScanConfigs, error) {
	var ret = models.ScanConfigs{
		Items: &[]models.ScanConfig{},
	}

	for _, config := range configs {
		sc, err := ConvertScanConfig(config)
		if err != nil {
			return nil, err
		}
		*ret.Items = append(*ret.Items, *sc)
	}
	ret.Total = int(total)

	return &ret, nil
}

func ConvertTarget(target *database.Target) (*models.Target, error) {
	var ret = models.Target{
		TargetInfo: &models.TargetType{},
	}

	switch target.Type {
	case "VMInfo":
		cloudProvider := models.CloudProvider(target.InstanceProvider)
		if err := ret.TargetInfo.FromVMInfo(models.VMInfo{
			InstanceID:       utils.StringPtr(target.InstanceID),
			InstanceProvider: &cloudProvider,
			Location:         utils.StringPtr(target.Location),
		}); err != nil {
			return nil, err
		}
	case "Dir":
		return nil, fmt.Errorf("unsupported target type Dir")
	case "Pod":
		return nil, fmt.Errorf("unsupported target type Pod")
	default:
		return nil, fmt.Errorf("unknown target type: %v", target.Type)
	}
	ret.Id = utils.StringPtr(target.ID.String())

	return &ret, nil
}

func ConvertTargets(targets []*database.Target, total int64) (*models.Targets, error) {
	var ret = models.Targets{
		Items: &[]models.Target{},
	}

	for _, target := range targets {
		tr, err := ConvertTarget(target)
		if err != nil {
			return nil, err
		}
		*ret.Items = append(*ret.Items, *tr)
	}

	ret.Total = int(total)

	return &ret, nil
}

func ConvertScanResult(scanResult *database.ScanResult) (*models.TargetScanResult, error) {
	var ret models.TargetScanResult

	if scanResult.Secrets != nil {
		ret.Secrets = &models.SecretScan{}
		if err := json.Unmarshal(scanResult.Secrets, ret.Secrets); err != nil {
			return nil, err
		}
	}
	if scanResult.Vulnerabilities != nil {
		ret.Vulnerabilities = &models.VulnerabilityScan{}
		if err := json.Unmarshal(scanResult.Vulnerabilities, ret.Vulnerabilities); err != nil {
			return nil, err
		}
	}

	if scanResult.Exploits != nil {
		ret.Exploits = &models.ExploitScan{}
		if err := json.Unmarshal(scanResult.Exploits, ret.Exploits); err != nil {
			return nil, err
		}
	}
	if scanResult.Malware != nil {
		ret.Malware = &models.MalwareScan{}
		if err := json.Unmarshal(scanResult.Malware, ret.Malware); err != nil {
			return nil, err
		}
	}
	if scanResult.Misconfigurations != nil {
		ret.Misconfigurations = &models.MisconfigurationScan{}
		if err := json.Unmarshal(scanResult.Misconfigurations, ret.Misconfigurations); err != nil {
			return nil, err
		}
	}
	if scanResult.Rootkits != nil {
		ret.Rootkits = &models.RootkitScan{}
		if err := json.Unmarshal(scanResult.Rootkits, ret.Rootkits); err != nil {
			return nil, err
		}
	}
	if scanResult.Sboms != nil {
		ret.Sboms = &models.SbomScan{}
		if err := json.Unmarshal(scanResult.Sboms, ret.Sboms); err != nil {
			return nil, err
		}
	}
	if scanResult.Status != nil {
		ret.Status = &models.TargetScanStatus{}
		if err := json.Unmarshal(scanResult.Status, ret.Status); err != nil {
			return nil, err
		}
	}
	ret.Id = utils.StringPtr(scanResult.ID.String())
	ret.ScanId = scanResult.ScanID
	ret.TargetId = scanResult.TargetID

	return &ret, nil
}

func ConvertScanResults(scanResults []*database.ScanResult, total int64) (*models.TargetScanResults, error) {
	var ret = models.TargetScanResults{
		Items: &[]models.TargetScanResult{},
	}

	for _, scanResult := range scanResults {
		sr, err := ConvertScanResult(scanResult)
		if err != nil {
			return nil, err
		}
		*ret.Items = append(*ret.Items, *sr)
	}

	ret.Total = int(total)

	return &ret, nil
}

func ConvertScan(scan *database.Scan) (*models.Scan, error) {
	var ret models.Scan

	if scan.ScanFamiliesConfig != nil {
		ret.ScanFamiliesConfig = &models.ScanFamiliesConfig{}
		if err := json.Unmarshal(scan.ScanFamiliesConfig, ret.ScanFamiliesConfig); err != nil {
			return nil, err
		}
	}

	if err := json.Unmarshal(scan.TargetIDs, ret.TargetIDs); err != nil {
		return nil, err
	}

	ret.Id = utils.StringPtr(scan.ID.String())
	ret.StartTime = scan.ScanStartTime
	ret.EndTime = scan.ScanEndTime
	ret.ScanConfigId = scan.ScanConfigId

	return &ret, nil
}

func ConvertScans(scans []*database.Scan, total int64) (*models.Scans, error) {
	var ret = models.Scans{
		Items: &[]models.Scan{},
	}

	for _, scan := range scans {
		sc, err := ConvertScan(scan)
		if err != nil {
			return nil, err
		}
		*ret.Items = append(*ret.Items, *sc)
	}

	ret.Total = int(total)

	return &ret, nil
}
