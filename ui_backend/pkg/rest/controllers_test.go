package rest

import (
	"reflect"
	"testing"

	"gotest.tools/v3/assert"

	backendmodels "github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/shared/pkg/utils"
	"github.com/openclarity/vmclarity/ui_backend/api/models"
)

func Test_getTargetLocation(t *testing.T) {
	targetInfo := backendmodels.TargetType{}
	err := targetInfo.FromVMInfo(backendmodels.VMInfo{
		Location: "us-east-1",
	})
	assert.NilError(t, err)

	type args struct {
		target backendmodels.Target
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "",
			args: args{
				target: backendmodels.Target{
					TargetInfo: &targetInfo,
				},
			},
			want:    "us-east-1",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getTargetLocation(tt.args.target)
			if (err != nil) != tt.wantErr {
				t.Errorf("getTargetLocation() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getTargetLocation() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_addTargetFindingsCount(t *testing.T) {
	type args struct {
		findingsCount map[backendmodels.ScanType]int
		summary       *backendmodels.ScanFindingsSummary
	}
	tests := []struct {
		name              string
		args              args
		wantFindingsCount map[backendmodels.ScanType]int
	}{
		{
			name: "from empty findings count",
			args: args{
				findingsCount: map[backendmodels.ScanType]int{},
				summary: &backendmodels.ScanFindingsSummary{
					TotalExploits:          utils.PointerTo(2),
					TotalMalware:           utils.PointerTo(3),
					TotalMisconfigurations: utils.PointerTo(4),
					TotalPackages:          utils.PointerTo(5),
					TotalRootkits:          utils.PointerTo(6),
					TotalSecrets:           utils.PointerTo(7),
					TotalVulnerabilities: &backendmodels.VulnerabilityScanSummary{
						TotalCriticalVulnerabilities:   utils.PointerTo(1),
						TotalHighVulnerabilities:       utils.PointerTo(2),
						TotalLowVulnerabilities:        utils.PointerTo(3),
						TotalMediumVulnerabilities:     utils.PointerTo(4),
						TotalNegligibleVulnerabilities: utils.PointerTo(5),
					},
				},
			},
			wantFindingsCount: map[backendmodels.ScanType]int{
				"Exploits":          2,
				"Malware":           3,
				"Misconfigurations": 4,
				"Rootkits":          6,
				"Secrets":           7,
				"Vulnerabilities":   15,
			},
		},
		{
			name: "from empty findings - only exploits",
			args: args{
				findingsCount: map[backendmodels.ScanType]int{},
				summary: &backendmodels.ScanFindingsSummary{
					TotalExploits: utils.PointerTo(2),
				},
			},
			wantFindingsCount: map[backendmodels.ScanType]int{
				"Exploits": 2,
			},
		},
		{
			name: "from existing findings - only exploits in summary",
			args: args{
				findingsCount: map[backendmodels.ScanType]int{
					"Exploits":          2,
					"Malware":           3,
					"Misconfigurations": 4,
					"Rootkits":          6,
					"Secrets":           7,
					"Vulnerabilities":   15,
				},
				summary: &backendmodels.ScanFindingsSummary{
					TotalExploits: utils.PointerTo(2),
				},
			},
			wantFindingsCount: map[backendmodels.ScanType]int{
				"Exploits":          4,
				"Malware":           3,
				"Misconfigurations": 4,
				"Rootkits":          6,
				"Secrets":           7,
				"Vulnerabilities":   15,
			},
		},
		{
			name: "from existing findings - add to all",
			args: args{
				findingsCount: map[backendmodels.ScanType]int{
					"Exploits":          2,
					"Malware":           3,
					"Misconfigurations": 4,
					"Rootkits":          6,
					"Secrets":           7,
					"Vulnerabilities":   15,
				},
				summary: &backendmodels.ScanFindingsSummary{
					TotalExploits:          utils.PointerTo(2),
					TotalMalware:           utils.PointerTo(3),
					TotalMisconfigurations: utils.PointerTo(4),
					TotalPackages:          utils.PointerTo(5),
					TotalRootkits:          utils.PointerTo(6),
					TotalSecrets:           utils.PointerTo(7),
					TotalVulnerabilities: &backendmodels.VulnerabilityScanSummary{
						TotalCriticalVulnerabilities:   utils.PointerTo(1),
						TotalHighVulnerabilities:       utils.PointerTo(2),
						TotalLowVulnerabilities:        utils.PointerTo(3),
						TotalMediumVulnerabilities:     utils.PointerTo(4),
						TotalNegligibleVulnerabilities: utils.PointerTo(5),
					},
				},
			},
			wantFindingsCount: map[backendmodels.ScanType]int{
				"Exploits":          4,
				"Malware":           6,
				"Misconfigurations": 8,
				"Rootkits":          12,
				"Secrets":           14,
				"Vulnerabilities":   30,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addTargetFindingsCount(tt.args.findingsCount, tt.args.summary)
			assert.DeepEqual(t, tt.args.findingsCount, tt.wantFindingsCount)
		})
	}
}

func Test_createFindingsCount(t *testing.T) {
	type args struct {
		findings map[backendmodels.ScanType]int
	}
	tests := []struct {
		name string
		args args
		want *models.FindingsCount
	}{
		{
			name: "all values exists",
			args: args{
				findings: map[backendmodels.ScanType]int{
					"Exploits":          2,
					"Malware":           3,
					"Misconfigurations": 4,
					"Rootkits":          6,
					"Secrets":           7,
					"Vulnerabilities":   15,
				},
			},
			want: &models.FindingsCount{
				Exploits:          utils.PointerTo(2),
				Malware:           utils.PointerTo(3),
				Misconfigurations: utils.PointerTo(4),
				Rootkits:          utils.PointerTo(6),
				Secrets:           utils.PointerTo(7),
				Vulnerabilities:   utils.PointerTo(15),
			},
		},
		{
			name: "not all values exists",
			args: args{
				findings: map[backendmodels.ScanType]int{
					"Exploits":          2,
					"Misconfigurations": 4,
					"Vulnerabilities":   15,
				},
			},
			want: &models.FindingsCount{
				Exploits:          utils.PointerTo(2),
				Malware:           utils.PointerTo(0),
				Misconfigurations: utils.PointerTo(4),
				Rootkits:          utils.PointerTo(0),
				Secrets:           utils.PointerTo(0),
				Vulnerabilities:   utils.PointerTo(15),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createFindingsCount(tt.args.findings); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createFindingsCount() = %v, want %v", got, tt.want)
			}
		})
	}
}
