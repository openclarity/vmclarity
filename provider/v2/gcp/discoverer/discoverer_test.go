package discoverer

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	apitypes "github.com/openclarity/vmclarity/api/types"
	"reflect"
)

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
			want: []string{},
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
			got := getZonesLastPart(tt.args.zones)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("getZonesLastPart() mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func Test_convertLabelsToTags(t *testing.T) {
	tests := []struct {
		name string
		args map[string]string
		want []apitypes.Tag
	}{
		{
			name: "sanity",
			args: map[string]string{
				"valid-tag": "valid-value",
			},
			want: []apitypes.Tag{{
				Key: "valid-tag", Value: "valid-value",
			}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertLabelsToTags(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertLabelsToTags() = %v, want %v", got, tt.want)
			}
		})
	}
}
