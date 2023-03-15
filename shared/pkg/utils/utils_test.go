package utils

import (
	"reflect"
	"testing"

	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/runtime_scan/pkg/utils"
)

func TestGetVulnerabilityTotalsPerSeverity(t *testing.T) {
	type args struct {
		vulnerabilities *[]models.Vulnerability
	}
	tests := []struct {
		name string
		args args
		want *models.VulnerabilityScanSummary
	}{
		{
			name: "nil should result in empty",
			args: args{
				vulnerabilities: nil,
			},
			want: &models.VulnerabilityScanSummary{
				TotalCriticalVulnerabilities:   utils.PointerTo[int](0),
				TotalHighVulnerabilities:       utils.PointerTo[int](0),
				TotalMediumVulnerabilities:     utils.PointerTo[int](0),
				TotalLowVulnerabilities:        utils.PointerTo[int](0),
				TotalNegligibleVulnerabilities: utils.PointerTo[int](0),
			},
		},
		{
			name: "check one type",
			args: args{
				vulnerabilities: utils.PointerTo[[]models.Vulnerability]([]models.Vulnerability{
					{
						Id: utils.PointerTo[string]("id1"),
						VulnerabilityInfo: &models.VulnerabilityInfo{
							Description:       utils.StringPtr("desc1"),
							Severity:          utils.PointerTo[models.VulnerabilitySeverity](models.CRITICAL),
							VulnerabilityName: utils.StringPtr("CVE-1"),
						},
					},
				}),
			},
			want: &models.VulnerabilityScanSummary{
				TotalCriticalVulnerabilities:   utils.PointerTo[int](1),
				TotalHighVulnerabilities:       utils.PointerTo[int](0),
				TotalMediumVulnerabilities:     utils.PointerTo[int](0),
				TotalLowVulnerabilities:        utils.PointerTo[int](0),
				TotalNegligibleVulnerabilities: utils.PointerTo[int](0),
			},
		},
		{
			name: "check all severity types",
			args: args{
				vulnerabilities: utils.PointerTo[[]models.Vulnerability]([]models.Vulnerability{
					{
						Id: utils.PointerTo[string]("id1"),
						VulnerabilityInfo: &models.VulnerabilityInfo{
							Description:       utils.StringPtr("desc1"),
							Severity:          utils.PointerTo[models.VulnerabilitySeverity](models.CRITICAL),
							VulnerabilityName: utils.StringPtr("CVE-1"),
						},
					},
					{
						Id: utils.PointerTo[string]("id2"),
						VulnerabilityInfo: &models.VulnerabilityInfo{
							Description:       utils.StringPtr("desc2"),
							Severity:          utils.PointerTo[models.VulnerabilitySeverity](models.HIGH),
							VulnerabilityName: utils.StringPtr("CVE-2"),
						},
					},
					{
						Id: utils.PointerTo[string]("id3"),
						VulnerabilityInfo: &models.VulnerabilityInfo{
							Description:       utils.StringPtr("desc3"),
							Severity:          utils.PointerTo[models.VulnerabilitySeverity](models.MEDIUM),
							VulnerabilityName: utils.StringPtr("CVE-3"),
						},
					},
					{
						Id: utils.PointerTo[string]("id4"),
						VulnerabilityInfo: &models.VulnerabilityInfo{
							Description:       utils.StringPtr("desc4"),
							Severity:          utils.PointerTo[models.VulnerabilitySeverity](models.LOW),
							VulnerabilityName: utils.StringPtr("CVE-4"),
						},
					},
					{
						Id: utils.PointerTo[string]("id5"),
						VulnerabilityInfo: &models.VulnerabilityInfo{
							Description:       utils.StringPtr("desc5"),
							Severity:          utils.PointerTo[models.VulnerabilitySeverity](models.NEGLIGIBLE),
							VulnerabilityName: utils.StringPtr("CVE-5"),
						},
					},
				}),
			},
			want: &models.VulnerabilityScanSummary{
				TotalCriticalVulnerabilities:   utils.PointerTo[int](1),
				TotalHighVulnerabilities:       utils.PointerTo[int](1),
				TotalMediumVulnerabilities:     utils.PointerTo[int](1),
				TotalLowVulnerabilities:        utils.PointerTo[int](1),
				TotalNegligibleVulnerabilities: utils.PointerTo[int](1),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetVulnerabilityTotalsPerSeverity(tt.args.vulnerabilities); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetVulnerabilityTotalsPerSeverity() = %v, want %v", got, tt.want)
			}
		})
	}
}
