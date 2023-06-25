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

func Test_isInScopeByTags(t *testing.T) {
	type args struct {
		vm          *computepb.Instance
		includeTags *[]models.Tag
		excludeTags *[]models.Tag
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "all tags are nil - in scope",
			args: args{
				vm: &computepb.Instance{
					Tags: nil,
				},
				includeTags: nil,
				excludeTags: nil,
			},
			want: true,
		},
		{
			name: "vm tags are not nil, user tags are nil - in scope",
			args: args{
				vm: &computepb.Instance{
					Tags: &computepb.Tags{
						Items: []string{"foo"},
					},
				},
				includeTags: nil,
				excludeTags: nil,
			},
			want: true,
		},
		{
			name: "vm tags are nil, include tags not - not in scope",
			args: args{
				vm: &computepb.Instance{
					Tags: nil,
				},
				includeTags: &[]models.Tag{
					{
						Key:   "foo",
						Value: "bar",
					},
				},
				excludeTags: nil,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isInScopeByTags(tt.args.vm, tt.args.includeTags, tt.args.excludeTags); got != tt.want {
				t.Errorf("isInScopeByTags() = %v, want %v", got, tt.want)
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

func Test_hasExcludeTags(t *testing.T) {
	type args struct {
		vm   *computepb.Instance
		tags *[]models.Tag
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "instance tags are nil",
			args: args{
				vm: &computepb.Instance{
					Tags: &computepb.Tags{},
				},
				tags: &[]models.Tag{
					{
						Key:   "key1",
						Value: "val1",
					},
					{
						Key:   "key2",
						Value: "val2",
					},
				},
			},
			want: false,
		},
		{
			name: "user tags are nil",
			args: args{
				vm: &computepb.Instance{
					Tags: &computepb.Tags{
						Items: []string{"key1=val1", "key2=val2"},
					},
				},
				tags: nil,
			},
			want: false,
		},
		{
			name: "has exclude tags",
			args: args{
				vm: &computepb.Instance{
					Tags: &computepb.Tags{
						Items: []string{"key1=val1", "key2=val2"},
					},
				},
				tags: &[]models.Tag{
					{
						Key:   "key1",
						Value: "val1",
					},
					{
						Key:   "key2",
						Value: "val2",
					},
				},
			},
			want: true,
		},
		{
			name: "does not have exclude tags",
			args: args{
				vm: &computepb.Instance{
					Tags: &computepb.Tags{
						Items: []string{"key3=val3", "key2=val2"},
					},
				},
				tags: &[]models.Tag{
					{
						Key:   "key1",
						Value: "val1",
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasExcludeTags(tt.args.vm, tt.args.tags); got != tt.want {
				t.Errorf("hasExcludeTags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_hasIncludeTags(t *testing.T) {
	type args struct {
		vm   *computepb.Instance
		tags *[]models.Tag
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "user include tags are nil",
			args: args{
				vm: &computepb.Instance{
					Tags: &computepb.Tags{
						Items: []string{"key1=val1", "key2=val2"},
					},
				},
				tags: nil,
			},
			want: true,
		},
		{
			name: "instance include tags are nil",
			args: args{
				vm: &computepb.Instance{
					Tags: &computepb.Tags{},
				},
				tags: &[]models.Tag{
					{
						Key:   "key1",
						Value: "val1",
					},
				},
			},
			want: false,
		},
		{
			name: "has include tags",
			args: args{
				vm: &computepb.Instance{
					Tags: &computepb.Tags{
						Items: []string{"key1=val1", "key2=val2"},
					},
				},
				tags: &[]models.Tag{
					{
						Key:   "key1",
						Value: "val1",
					},
				},
			},
			want: true,
		},
		{
			name: "does not have include tags",
			args: args{
				vm: &computepb.Instance{
					Tags: &computepb.Tags{
						Items: []string{"key3=val3", "key2=val2"},
					},
				},
				tags: &[]models.Tag{
					{
						Key:   "key1",
						Value: "val1",
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasIncludeTags(tt.args.vm, tt.args.tags); got != tt.want {
				t.Errorf("hasIncludeTags() = %v, want %v", got, tt.want)
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
