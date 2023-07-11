package gcp

import (
	"reflect"
	"testing"

	"cloud.google.com/go/compute/apiv1/computepb"

	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/shared/pkg/utils"
)

func Test_convertTags(t *testing.T) {
	type args struct {
		tags *computepb.Tags
	}
	tests := []struct {
		name string
		args args
		want *[]models.Tag
	}{
		{
			name: "sanity",
			args: args{
				tags: &computepb.Tags{
					Items: []string{"tag1", "tag2=val2", "tag3="},
				},
			},
			want: &[]models.Tag{
				{
					Key:   "tag1",
					Value: "",
				},
				{
					Key:   "tag2",
					Value: "val2",
				},
				{
					Key:   "tag3",
					Value: "",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertTags(tt.args.tags); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertTags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_convertTagsToMap(t *testing.T) {
	type args struct {
		tags *computepb.Tags
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "Items=nil",
			args: args{
				tags: &computepb.Tags{},
			},
			want: map[string]string{},
		},
		{
			name: "no tags",
			args: args{
				tags: &computepb.Tags{
					Items: []string{},
				},
			},
			want: map[string]string{},
		},
		{
			name: "sanity",
			args: args{
				tags: &computepb.Tags{
					Items: []string{"key1=val1", "key2"},
				},
			},
			want: map[string]string{
				"key1": "val1",
				"key2": "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertTagsToMap(tt.args.tags); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertTagsToMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
							DeviceName: utils.PointerTo("device1"),
							Boot:       utils.PointerTo(true),
						},
						{
							DeviceName: utils.PointerTo("device2"),
							Boot:       utils.PointerTo(false),
						},
					},
				},
			},
			want: &computepb.AttachedDisk{
				DeviceName: utils.PointerTo("device1"),
				Boot:       utils.PointerTo(true),
			},
			wantErr: false,
		},
		{
			name: "not found",
			args: args{
				vm: &computepb.Instance{
					Disks: []*computepb.AttachedDisk{
						{
							DeviceName: utils.PointerTo("device1"),
							Boot:       utils.PointerTo(false),
						},
						{
							DeviceName: utils.PointerTo("device2"),
							Boot:       utils.PointerTo(false),
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

func Test_getZonesLastPart(t *testing.T) {
	type args struct {
		zones []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "empty",
			args: args{
				zones: []string{},
			},
			want: nil,
		},
		{
			name: "get two zones",
			args: args{
				zones: []string{"https://www.googleapis.com/compute/v1/projects/gcp-etigcp-nprd-12855/zones/us-central1-c", "https://www.googleapis.com/compute/v1/projects/gcp-etigcp-nprd-12855/zones/us-central1-a"},
			},
			want: []string{"us-central1-c", "us-central1-a"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getZonesLastPart(tt.args.zones); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getZonesLastPart() = %v, want %v", got, tt.want)
			}
		})
	}
}
