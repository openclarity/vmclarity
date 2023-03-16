package gorm

import (
	"testing"
	"time"

	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/shared/pkg/utils"
)

func Test_isEmptyOperationTime(t *testing.T) {
	now := time.Now()
	type args struct {
		operationTime *time.Time
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "nil operationTime",
			args: args{
				operationTime: nil,
			},
			want: true,
		},
		{
			name: "zero operationTime",
			args: args{
				operationTime: &time.Time{},
			},
			want: true,
		},
		{
			name: "not empty operationTime",
			args: args{
				operationTime: &now,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isEmptyOperationTime(tt.args.operationTime); got != tt.want {
				t.Errorf("isEmptyOperationTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_validateRuntimeScheduleScanConfig(t *testing.T) {
	type args struct {
		scheduled *models.RuntimeScheduleScanConfig
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nil scheduled",
			args: args{
				scheduled: nil,
			},
			wantErr: true,
		},
		{
			name: "both cron and operation time is missing",
			args: args{
				scheduled: &models.RuntimeScheduleScanConfig{
					CronLine:      nil,
					OperationTime: nil,
				},
			},
			wantErr: true,
		},
		{
			name: "both cron and operation time is missing, time is empty",
			args: args{
				scheduled: &models.RuntimeScheduleScanConfig{
					CronLine:      nil,
					OperationTime: &time.Time{},
				},
			},
			wantErr: true,
		},
		{
			name: "operation time is missing - not a valid cron expression",
			args: args{
				scheduled: &models.RuntimeScheduleScanConfig{
					CronLine: utils.PointerTo("not valid"),
				},
			},
			wantErr: true,
		},
		{
			name: "operation time is missing - operation time should be set",
			args: args{
				scheduled: &models.RuntimeScheduleScanConfig{
					CronLine: utils.PointerTo("0 */4 * * *"),
				},
			},
			wantErr: false,
		},
		{
			name: "cron line is missing - do nothing",
			args: args{
				scheduled: &models.RuntimeScheduleScanConfig{
					OperationTime: utils.PointerTo(time.Now()),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateRuntimeScheduleScanConfig(tt.args.scheduled)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateRuntimeScheduleScanConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && tt.args.scheduled.OperationTime == nil {
				t.Errorf("validateRuntimeScheduleScanConfig() operation time must be set after successfull validation")
			}
		})
	}
}
