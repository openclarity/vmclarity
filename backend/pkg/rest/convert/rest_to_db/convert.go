package rest_to_db

import (
	"encoding/json"
	"fmt"

	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/backend/pkg/database"
)

func ConvertScanConfig(config *models.ScanConfig) (*database.ScanConfig, error) {
	var ret database.ScanConfig
	var err error

	ret.ScanFamiliesConfig, err = json.Marshal(config.ScanFamiliesConfig)
	if err != nil {
		return nil, err
	}
	ret.Scope, err = config.Scope.MarshalJSON()
	if err != nil {
		return nil, err
	}
	ret.Scheduled, err = config.Scheduled.MarshalJSON()
	if err != nil {
		return nil, err
	}

	ret.Name = *config.Name

	return &ret, nil
}

func ConvertTarget(target *models.Target) (*database.Target, error) {
	disc, err := target.TargetInfo.Discriminator()
	if err != nil {
		return nil, err
	}
	switch disc {
	case "VMInfo":
		vminfo, err := target.TargetInfo.AsVMInfo()
		if err != nil {
			return nil, err
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

func ConvertScanResult(result *models.TargetScanResult) (*database.ScanResult, error) {
	var ret database.ScanResult
	var err error

	ret.ScanID = result.ScanId
	ret.TargetID = result.TargetId

	ret.Exploits, err = json.Marshal(result.Exploits)
	if err != nil {
		return nil, err
	}
	ret.Malware, err = json.Marshal(result.Malware)
	if err != nil {
		return nil, err
	}
	ret.Misconfigurations, err = json.Marshal(result.Misconfigurations)
	if err != nil {
		return nil, err
	}
	ret.Rootkits, err = json.Marshal(result.Rootkits)
	if err != nil {
		return nil, err
	}
	ret.Sboms, err = json.Marshal(result.Sboms)
	if err != nil {
		return nil, err
	}
	ret.Secrets, err = json.Marshal(result.Secrets)
	if err != nil {
		return nil, err
	}
	ret.Status, err = json.Marshal(result.Status)
	if err != nil {
		return nil, err
	}
	ret.Vulnerabilities, err = json.Marshal(result.Vulnerabilities)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}

func ConvertScan(scan *models.Scan) (*database.Scan, error) {
	var ret database.Scan
	var err error

	ret.ScanConfigId = *scan.ScanConfigId
	ret.ScanEndTime = *scan.EndTime
	ret.ScanStartTime = *scan.StartTime
	ret.ScanFamiliesConfig, err = json.Marshal(scan.ScanFamiliesConfig)
	if err != nil {
		return nil, err
	}
	ret.TargetIDs, err = json.Marshal(scan.TargetIDs)
	if err != nil {
		return nil, err
	}

	return &ret, nil
}
