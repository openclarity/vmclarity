package misconfiguration

import (
	"reflect"
	"testing"

	"github.com/openclarity/vmclarity/shared/pkg/families/misconfiguration/types"
)

func TestStripPathFromResult(t *testing.T) {
	type args struct {
		result types.ScannerResult
		path   string
	}
	tests := []struct {
		name string
		args args
		want types.ScannerResult
	}{
		{
			name: "sanity",
			args: args{
				result: types.ScannerResult{
					ScannerName: "scanner1",
					Misconfigurations: []types.Misconfiguration{
						{
							ScannedPath:     "/mnt/foo",
							TestCategory:    "test1",
							TestID:          "id1",
							TestDescription: "desc1",
						},
						{
							ScannedPath:     "/mnt/foo2",
							TestCategory:    "test2",
							TestID:          "id2",
							TestDescription: "desc2",
						},
						{
							ScannedPath:     "/foo3",
							TestCategory:    "test3",
							TestID:          "id3",
							TestDescription: "desc3",
						},
					},
				},
				path: "/mnt",
			},
			want: types.ScannerResult{
				ScannerName: "scanner1",
				Misconfigurations: []types.Misconfiguration{
					{
						ScannedPath:     "/foo",
						TestCategory:    "test1",
						TestID:          "id1",
						TestDescription: "desc1",
					},
					{
						ScannedPath:     "/foo2",
						TestCategory:    "test2",
						TestID:          "id2",
						TestDescription: "desc2",
					},
					{
						ScannedPath:     "/foo3",
						TestCategory:    "test3",
						TestID:          "id3",
						TestDescription: "desc3",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StripPathFromResult(tt.args.result, tt.args.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StripPathFromResult() = %v, want %v", got, tt.want)
			}
		})
	}
}
