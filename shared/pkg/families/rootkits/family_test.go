package rootkits

import (
	"reflect"
	"testing"

	"github.com/openclarity/vmclarity/shared/pkg/families/rootkits/common"
)

func TestStripPathFromResult(t *testing.T) {
	type args struct {
		result *common.Results
		path   string
	}
	tests := []struct {
		name string
		args args
		want *common.Results
	}{
		{
			name: "sanity",
			args: args{
				result: &common.Results{
					Rootkits: []common.Rootkit{
						{
							Message:     "rootkit found in /mnt/foo path",
							RootkitName: "rk1",
							RootkitType: "t1",
						},
						{
							Message:     "rootkit found in /mnt/bar path",
							RootkitName: "rk2",
							RootkitType: "t2",
						},
					},
					ScannedInput: "/mnt/foo",
					ScannerName:  "scanner1",
				},
				path: "/mnt",
			},
			want: &common.Results{
				Rootkits: []common.Rootkit{
					{
						Message:     "rootkit found in /foo path",
						RootkitName: "rk1",
						RootkitType: "t1",
					},
					{
						Message:     "rootkit found in /bar path",
						RootkitName: "rk2",
						RootkitType: "t2",
					},
				},
				ScannedInput: "/mnt/foo",
				ScannerName:  "scanner1",
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
