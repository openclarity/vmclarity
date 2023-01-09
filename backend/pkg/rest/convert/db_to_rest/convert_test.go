package db_to_rest

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
		config *database.ScanConfig
	}
	tests := []struct {
		name    string
		args    args
		want    *models.ScanConfig
		wantErr bool
	}{
		{
			name: "",
			args: args{
				config: &database.ScanConfig{
					Model: gorm.Model{
						ID: 1,
					},
					Name:               "test",
					ScanFamiliesConfig: scanFamiliesConfigB,
					Scheduled:          runtimeScheduleScanConfigTypeB,
					Scope:              scanScopeTypeB,
				},
			},
			want: &models.ScanConfig{
				Id:                 utils.StringPtr("1"),
				Name:               utils.StringPtr("test"),
				ScanFamiliesConfig: &scanFamiliesConfig,
				Scheduled:          &runtimeScheduleScanConfigType,
				Scope:              &scanScopeType,
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

func TestConvertScanConfigs(t *testing.T) {
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
		configs []*database.ScanConfig
		total   int64
	}
	tests := []struct {
		name    string
		args    args
		want    *models.ScanConfigs
		wantErr bool
	}{
		{
			name: "",
			args: args{
				configs: []*database.ScanConfig{
					{
						Model: gorm.Model{
							ID: 1,
						},
						Name:               "test",
						ScanFamiliesConfig: scanFamiliesConfigB,
						Scheduled:          runtimeScheduleScanConfigTypeB,
						Scope:              scanScopeTypeB,
					},
				},
				total: 1,
			},
			want: &models.ScanConfigs{
				Items: &[]models.ScanConfig{
					{
						Id:                 utils.StringPtr("1"),
						Name:               utils.StringPtr("test"),
						ScanFamiliesConfig: &scanFamiliesConfig,
						Scheduled:          &runtimeScheduleScanConfigType,
						Scope:              &scanScopeType,
					},
				},
				Total: 1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertScanConfigs(tt.args.configs, tt.args.total)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertScanConfigs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConvertScanConfigs() got = %v, want %v", got, tt.want)
			}
		})
	}
}
