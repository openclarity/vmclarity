package findingkey

import (
	"reflect"
	"testing"

	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/pkg/shared/utils"
)

func TestGenerateInfoFinderKey(t *testing.T) {
	type args struct {
		info models.InfoFinderFindingInfo
	}
	tests := []struct {
		name string
		args args
		want InfoFinderKey
	}{
		{
			name: "sanity",
			args: args{
				info: models.InfoFinderFindingInfo{
					Data:        utils.PointerTo("data"),
					Path:        utils.PointerTo("path"),
					ScannerName: utils.PointerTo("scanner"),
					Type:        utils.PointerTo(models.InfoTypeSSHAuthorizedKeyFingerprint),
				},
			},
			want: InfoFinderKey{
				ScannerName: "scanner",
				Type:        string(models.InfoTypeSSHAuthorizedKeyFingerprint),
				Data:        "data",
				Path:        "path",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateInfoFinderKey(tt.args.info); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenerateInfoFinderKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInfoFinderKey_String(t *testing.T) {
	type fields struct {
		ScannerName string
		Type        string
		Data        string
		Path        string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "sanity",
			fields: fields{
				ScannerName: "scanner",
				Type:        string(models.InfoTypeSSHAuthorizedKeyFingerprint),
				Data:        "data",
				Path:        "path",
			},
			want: "scanner.SSHAuthorizedKeyFingerprint.data.path",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := InfoFinderKey{
				ScannerName: tt.fields.ScannerName,
				Type:        tt.fields.Type,
				Data:        tt.fields.Data,
				Path:        tt.fields.Path,
			}
			if got := k.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
