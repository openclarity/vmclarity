package rest_to_db

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/backend/pkg/database"
	"github.com/openclarity/vmclarity/runtime_scan/pkg/utils"
	"gorm.io/gorm"
	"gotest.tools/v3/assert"
)

func TestConvertScanConfig(t *testing.T) {
	scanFamiliesConfig := models.ScanFamiliesConfig{
		Vulnerabilities: &models.VulnerabilitiesConfig{Enabled: utils.BoolPtr(true)},
	}

	scanFamiliesConfigB, err := json.Marshal(&scanFamiliesConfig)
	assert.NilError(t, err)

	awsScanScope := models.AwsScanScope{
		All:                        utils.BoolPtr(true),
		InstanceTagExclusion:       nil,
		InstanceTagSelector:        nil,
		ObjectType:                 "AwsScanScope",
		Regions:                    nil,
		ShouldScanStoppedInstances: utils.BoolPtr(false),
	}

	var scanScopeType models.ScanScopeType

	err = scanScopeType.FromAwsScanScope(awsScanScope)
	assert.NilError(t, err)

	scanScopeTypeB, err := scanScopeType.MarshalJSON()
	assert.NilError(t, err)

	var byHoursScheduleScanConfig = models.ByHoursScheduleScanConfig{
		HoursInterval: utils.IntPtr(2),
		ObjectType:    "ByHoursScheduleScanConfig",
	}

	var runtimeScheduleScanConfigType models.RuntimeScheduleScanConfigType
	err = runtimeScheduleScanConfigType.FromByHoursScheduleScanConfig(byHoursScheduleScanConfig)
	assert.NilError(t, err)

	runtimeScheduleScanConfigTypeB, err := runtimeScheduleScanConfigType.MarshalJSON()
	assert.NilError(t, err)

	type args struct {
		config *models.ScanConfig
	}
	tests := []struct {
		name    string
		args    args
		want    *database.ScanConfig
		wantErr bool
	}{
		{
			name: "",
			args: args{
				config: &models.ScanConfig{
					Id:                 utils.StringPtr("1"),
					Name:               utils.StringPtr("scanConfigName"),
					ScanFamiliesConfig: &scanFamiliesConfig,
					Scheduled:          &runtimeScheduleScanConfigType,
					Scope:              &scanScopeType,
				},
			},
			want: &database.ScanConfig{
				Name:               "scanConfigName",
				ScanFamiliesConfig: scanFamiliesConfigB,
				Scheduled:          runtimeScheduleScanConfigTypeB,
				Scope:              scanScopeTypeB,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertScanConfig(tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertScanConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertScanConfig() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertScanResult(t *testing.T) {
	vulnerabilities := []models.Vulnerability{
		{
			VulnerabilityInfo: &models.VulnerabilityInfo{
				Description:       utils.StringPtr("desc"),
				Id:                utils.StringPtr("1"),
				VulnerabilityName: utils.StringPtr("name"),
			},
			Id: utils.StringPtr("1"),
		},
	}

	vulnerabilitiesB, err := json.Marshal(vulnerabilities)
	assert.NilError(t, err)

	type args struct {
		result *models.TargetScanResult
	}
	tests := []struct {
		name    string
		args    args
		want    *database.ScanResult
		wantErr bool
	}{
		{
			name: "",
			args: args{
				result: &models.TargetScanResult{
					Id:       nil,
					ScanId:   "3",
					TargetId: "2",
					Vulnerabilities: &models.VulnerabilityScan{
						Vulnerabilities: &vulnerabilities,
					},
				},
			},
			want: &database.ScanResult{
				Model: gorm.Model{
					ID: 1,
				},
				ScanID:          "3",
				TargetID:        "2",
				Vulnerabilities: vulnerabilitiesB,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertScanResult(tt.args.result)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertScanResult() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertScanResult() got = %v, want %v", got, tt.want)
			}
		})
	}
}
