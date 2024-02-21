package scanner

import (
	"reflect"
	"testing"

	"cloud.google.com/go/compute/apiv1/computepb"

	"github.com/openclarity/vmclarity/core/to"
)

func Test_getInstanceBootDisk(t *testing.T) {
	type args struct {
		vm *computepb.Instance
	}
	tests := []struct {
		name    string
		args    args
		want    *computepb.AttachedDisk
		wantErr bool
	}{
		{
			name: "found",
			args: args{
				vm: &computepb.Instance{
					Disks: []*computepb.AttachedDisk{
						{
							DeviceName: to.Ptr("device1"),
							Boot:       to.Ptr(true),
						},
						{
							DeviceName: to.Ptr("device2"),
							Boot:       to.Ptr(false),
						},
					},
				},
			},
			want: &computepb.AttachedDisk{
				DeviceName: to.Ptr("device1"),
				Boot:       to.Ptr(true),
			},
			wantErr: false,
		},
		{
			name: "not found",
			args: args{
				vm: &computepb.Instance{
					Disks: []*computepb.AttachedDisk{
						{
							DeviceName: to.Ptr("device1"),
							Boot:       to.Ptr(false),
						},
						{
							DeviceName: to.Ptr("device2"),
							Boot:       to.Ptr(false),
						},
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getInstanceBootDisk(tt.args.vm)
			if (err != nil) != tt.wantErr {
				t.Errorf("getInstanceBootDisk() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getInstanceBootDisk() got = %v, want %v", got, tt.want)
			}
		})
	}
}
