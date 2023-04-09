package chkrootkit

import (
	"encoding/json"
	"os"
	"testing"

	"gotest.tools/v3/assert"

	chkrootkitutils "github.com/openclarity/vmclarity/shared/pkg/families/rootkits/chkrootkit/utils"
	"github.com/openclarity/vmclarity/shared/pkg/families/rootkits/common"
)

func Test_toResultsRootkits(t *testing.T) {
	chkrootkitOutput, err := os.ReadFile("utils/testdata/chkrootkit_output.txt")
	assert.NilError(t, err)
	rootkits, err := chkrootkitutils.ParseChkrootkitOutput(chkrootkitOutput)
	assert.NilError(t, err)

	type args struct {
		rootkits []chkrootkitutils.Rootkit
	}
	tests := []struct {
		name string
		args args
		want []common.Rootkit
	}{
		{
			name: "sanity",
			args: args{
				rootkits: rootkits,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toResultsRootkits(tt.args.rootkits)
			//if diff := cmp.Diff(tt.want, got); diff != "" {
			//	t.Errorf("toResultsRootkits() mismatch (-want +got):\n%s", diff)
			//}
			t.Logf("toResultsRootkits() results: %+v", prettyPrint(t, got))
		})
	}
}

func prettyPrint(t *testing.T, got any) string {
	jsonResults, err := json.MarshalIndent(got, "", "    ")
	assert.NilError(t, err)
	return string(jsonResults)
}
