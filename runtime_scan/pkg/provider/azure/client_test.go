package azure

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v4"
	"github.com/openclarity/vmclarity/shared/pkg/utils"
)

func Test_isEncrypted(t *testing.T) {
	type args struct {
		disk armcompute.DisksClientGetResponse
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "encrypted",
			args: args{
				disk: armcompute.DisksClientGetResponse{
					Disk: armcompute.Disk{
						Properties: &armcompute.DiskProperties{
							EncryptionSettingsCollection: &armcompute.EncryptionSettingsCollection{
								Enabled: utils.PointerTo(true),
							},
						},
					},
				},
			},
			want: true,
		},
		{
			name: "not encrypted",
			args: args{
				disk: armcompute.DisksClientGetResponse{
					Disk: armcompute.Disk{
						Properties: &armcompute.DiskProperties{
							EncryptionSettingsCollection: nil,
						},
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isEncrypted(tt.args.disk); got != tt.want {
				t.Errorf("isEncrypted() = %v, want %v", got, tt.want)
			}
		})
	}
}
