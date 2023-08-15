package scanestimationwatcher

import (
	"testing"

	"gotest.tools/v3/assert"

	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/pkg/shared/utils"
)

func Test_updateScanEstimationSummaryFromAssetScanEstimation(t *testing.T) {
	type args struct {
		scanEstimation *models.ScanEstimation
		result         models.AssetScanEstimation
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "",
			args: args{
				scanEstimation: &models.ScanEstimation{
					Summary: &models.ScanEstimationSummary{
						JobsCompleted: utils.PointerTo(0),
						JobsLeftToRun: utils.PointerTo(0),
					},
				},
				result: models.AssetScanEstimation{
					State: &models.AssetScanEstimationState{
						State: utils.PointerTo(models.AssetScanEstimationStateStateFailed),
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := updateScanEstimationSummaryFromAssetScanEstimation(tt.args.scanEstimation, tt.args.result); (err != nil) != tt.wantErr {
				t.Errorf("updateScanEstimationSummaryFromAssetScanEstimation() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Assert(t, *tt.args.scanEstimation.Summary.JobsLeftToRun == 0)
			assert.Assert(t, *tt.args.scanEstimation.Summary.JobsCompleted == 1)
		})
	}
}
