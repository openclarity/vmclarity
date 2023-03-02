package cmd

import (
	"reflect"
	"testing"

	"github.com/openclarity/kubeclarity/shared/pkg/utils"

	"github.com/openclarity/vmclarity/shared/pkg/families"
	"github.com/openclarity/vmclarity/shared/pkg/families/sbom"
	"github.com/openclarity/vmclarity/shared/pkg/families/secrets"
	"github.com/openclarity/vmclarity/shared/pkg/families/vulnerabilities"
)

func Test_isSupportedFS(t *testing.T) {
	type args struct {
		fs string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "supported ext4",
			args: args{
				fs: fsTypeExt4,
			},
			want: true,
		},
		{
			name: "supported xfs",
			args: args{
				fs: fsTypeXFS,
			},
			want: true,
		},
		{
			name: "not supported btrfs",
			args: args{
				fs: "btrfs",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isSupportedFS(tt.args.fs); got != tt.want {
				t.Errorf("isSupportedFS() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_setMountPointsForFamiliesInput(t *testing.T) {
	type args struct {
		mountPoints    []string
		familiesConfig *families.Config
	}
	tests := []struct {
		name string
		args args
		want *families.Config
	}{
		{
			name: "sbom, vuls and secrets are enabled",
			args: args{
				mountPoints: []string{"/mnt/snapshot1"},
				familiesConfig: &families.Config{
					SBOM: sbom.Config{
						Enabled: true,
						Inputs:  nil,
					},
					Vulnerabilities: vulnerabilities.Config{
						Enabled:       true,
						Inputs:        nil,
						InputFromSbom: false,
					},
					Secrets: secrets.Config{
						Enabled: true,
						Inputs:  nil,
					},
				},
			},
			want: &families.Config{
				SBOM: sbom.Config{
					Enabled: true,
					Inputs: []sbom.Input{
						{
							Input:     "/mnt/snapshot1",
							InputType: string(utils.ROOTFS),
						},
					},
				},
				Vulnerabilities: vulnerabilities.Config{
					Enabled:       true,
					InputFromSbom: true,
				},
				Secrets: secrets.Config{
					Enabled: true,
					Inputs: []secrets.Input{
						{
							Input:     "/mnt/snapshot1",
							InputType: string(utils.ROOTFS),
						},
					},
				},
			},
		},
		{
			name: "only vuls enabled",
			args: args{
				mountPoints: []string{"/mnt/snapshot1"},
				familiesConfig: &families.Config{
					Vulnerabilities: vulnerabilities.Config{
						Enabled:       true,
						Inputs:        nil,
						InputFromSbom: false,
					},
				},
			},
			want: &families.Config{
				Vulnerabilities: vulnerabilities.Config{
					Enabled: true,
					Inputs: []vulnerabilities.Input{
						{
							Input:     "/mnt/snapshot1",
							InputType: string(utils.ROOTFS),
						},
					},
					InputFromSbom: false,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := setMountPointsForFamiliesInput(tt.args.mountPoints, tt.args.familiesConfig); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("setMountPointsForFamiliesInput() = %v, want %v", got, tt.want)
			}
		})
	}
}
