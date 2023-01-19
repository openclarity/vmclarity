package configwatcher

import (
	"testing"
	"time"

	"github.com/openclarity/vmclarity/api/models"
)

func Test_hasRunningOrCompletedScan(t *testing.T) {
	testScanConfigID := "testID"
	otherScanConfigID := "otherID"
	operationTime := time.Now()
	afterOperationTime := operationTime.Add(time.Minute * 5)
	beforeOperationTime := operationTime.Add(-time.Minute * 5)
	type args struct {
		scans         *models.Scans
		scanConfigID  string
		operationTime time.Time
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "scans.Items is null",
			args: args{
				scans:         &models.Scans{},
				scanConfigID:  testScanConfigID,
				operationTime: operationTime,
			},
			want: false,
		},
		{
			name: "scans.Items is empty list",
			args: args{
				scans: &models.Scans{
					Items: &[]models.Scan{},
				},
				scanConfigID:  testScanConfigID,
				operationTime: operationTime,
			},
			want: false,
		},
		{
			name: "there are no scans with scan config ID",
			args: args{
				scans: &models.Scans{
					Items: &[]models.Scan{
						{
							ScanConfigId: &otherScanConfigID,
						},
					},
				},
				scanConfigID:  testScanConfigID,
				operationTime: operationTime,
			},
			want: false,
		},
		{
			name: "there is a scans without end time",
			args: args{
				scans: &models.Scans{
					Items: &[]models.Scan{
						{
							ScanConfigId: &testScanConfigID,
						},
					},
				},
				scanConfigID:  testScanConfigID,
				operationTime: operationTime,
			},
			want: true,
		},
		{
			name: "there is a scans with end time and start time after operation time",
			args: args{
				scans: &models.Scans{
					Items: &[]models.Scan{
						{
							ScanConfigId: &testScanConfigID,
							StartTime:    &afterOperationTime,
							EndTime:      &operationTime,
						},
					},
				},
				scanConfigID:  testScanConfigID,
				operationTime: operationTime,
			},
			want: true,
		},
		{
			name: "there is a scans with end time and start time before operation time",
			args: args{
				scans: &models.Scans{
					Items: &[]models.Scan{
						{
							ScanConfigId: &testScanConfigID,
							StartTime:    &beforeOperationTime,
							EndTime:      &operationTime,
						},
					},
				},
				scanConfigID:  testScanConfigID,
				operationTime: operationTime,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hasRunningOrCompletedScan(tt.args.scans, tt.args.scanConfigID, tt.args.operationTime); got != tt.want {
				t.Errorf("hasRunningOrCompletedScan() = %v, want %v", got, tt.want)
			}
		})
	}
}
