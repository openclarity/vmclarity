// Copyright © 2023 Cisco Systems, Inc. and its affiliates.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package scheduler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/openclarity/vmclarity/api/client"
	"github.com/openclarity/vmclarity/api/models"
	_config "github.com/openclarity/vmclarity/runtime_scan/pkg/config"
	"github.com/openclarity/vmclarity/runtime_scan/pkg/provider"
	_scanner "github.com/openclarity/vmclarity/runtime_scan/pkg/scanner"
	"github.com/openclarity/vmclarity/runtime_scan/pkg/types"
	"github.com/openclarity/vmclarity/runtime_scan/pkg/utils"
)

type Scheduler struct {
	stopChan        chan struct{}
	scanConfigChan  chan *map[string]models.ScanConfig
	scannerConfig   *_config.ScannerConfig
	providerClient  provider.Client
	backendClient   *client.ClientWithResponses
	scheduledCancel *map[string]context.CancelFunc
}

type Params struct {
	Interval   time.Duration
	StartTime  time.Time
	SingleScan bool
}

func CreateScheduler(scanConfigChan chan *map[string]models.ScanConfig,
	scannerConfig *_config.ScannerConfig,
	providerClient provider.Client,
	backendClient *client.ClientWithResponses,
) *Scheduler {
	return &Scheduler{
		stopChan:       make(chan struct{}),
		scanConfigChan: scanConfigChan,
		scannerConfig:  scannerConfig,
		providerClient: providerClient,
		backendClient:  backendClient,
	}
}

func (s *Scheduler) Start(errChan chan struct{}) {
	// Clear
	close(s.stopChan)
	s.stopChan = make(chan struct{})
	for {
		select {
		case scanConfigMap := <-s.scanConfigChan:
			if err := s.scheduleNewScans(scanConfigMap); err != nil {
				if errChan != nil {
					errChan <- struct{}{}
				}
			}
		case <-s.stopChan:
			log.Infof("Stop watching scan configs.")
			return
		}
	}
}

func (s *Scheduler) Stop() {
	s.stopChan <- struct{}{}
}

func (s *Scheduler) scheduleNewScans(scanConfigMap *map[string]models.ScanConfig) error {
	for _, scanConfig := range *scanConfigMap {
		ctx, cancel := context.WithCancel(context.Background())
		params, err := handleNewScheduleScanConfig(scanConfig.Scheduled)
		if err != nil {
			return fmt.Errorf("failed to schedule new scan with scanConfigID %s: %v", *scanConfig.Id, err)
		}

		s.schedule(ctx, params, scanConfig)
	}

	return nil
}

const (
	secondsInHour = 60 * 60
	secondsInDay  = 24 * secondsInHour
	secondsInWeek = 7 * secondsInDay
)

func getIntervalAndStartTimeFromByDaysScheduleScanConfig(timeNow time.Time, scanConfig *models.ByDaysScheduleScanConfig) (time.Duration, time.Time) {
	interval := time.Duration(*scanConfig.DaysInterval*secondsInDay) * time.Second
	hour := int(*scanConfig.TimeOfDay.Hour)
	minute := int(*scanConfig.TimeOfDay.Minute)
	year, month, day := timeNow.Date()

	startTime := time.Date(year, month, day, hour, minute, 0, 0, time.UTC)

	return interval, startTime
}

func getIntervalAndStartTimeFromByHoursScheduleScanConfig(timeNow time.Time, scanConfig *models.ByHoursScheduleScanConfig) (time.Duration, time.Time) {
	interval := time.Duration(*scanConfig.HoursInterval*secondsInHour) * time.Second

	return interval, timeNow
}

func getIntervalAndStartTimeFromWeeklyScheduleScanConfig(timeNow time.Time, scanConfig *models.WeeklyScheduleScanConfig) (time.Duration, time.Time) {
	interval := time.Duration(secondsInWeek) * time.Second

	currentDay := timeNow.Weekday() + 1
	diffDays := int64(*scanConfig.DayInWeek) - int64(currentDay)

	hour := int(*scanConfig.TimeOfDay.Hour)
	minute := int(*scanConfig.TimeOfDay.Minute)
	year, month, day := timeNow.Add(time.Duration(diffDays*secondsInDay) * time.Second).Date()

	startTime := time.Date(year, month, day, hour, minute, 0, 0, time.UTC)

	return interval, startTime
}

func handleNewScheduleScanConfig(scheduleScanConfigType *models.RuntimeScheduleScanConfigType) (*Params, error) {
	var interval time.Duration
	var startTime time.Time
	singleScan := false

	timeNow := time.Now().UTC()

	scanConfigType, err := scheduleScanConfigType.ValueByDiscriminator()
	if err != nil {
		return nil, fmt.Errorf("failed to determine scheduled scan config type: %v", err)
	}
	switch scanConfigType.(type) {
	case models.SingleScheduleScanConfig:
		var err error
		singleScan = true
		// nolint:forcetypeassert
		scanConfig := scanConfigType.(*models.SingleScheduleScanConfig)
		startTime, err = time.Parse(time.RFC3339, scanConfig.OperationTime.String())
		if err != nil {
			return nil, fmt.Errorf("failed to parse operation time: %v. %v", scanConfig.OperationTime.String(), err)
		}
		// set interval to a positive value, so we will not crash when starting ticker in Scheduler.spin. This will not be used.
		interval = 1
	case models.ByHoursScheduleScanConfig:
		// nolint:forcetypeassert
		scanConfig := scanConfigType.(*models.ByHoursScheduleScanConfig)
		interval, startTime = getIntervalAndStartTimeFromByHoursScheduleScanConfig(timeNow, scanConfig)
	case models.ByDaysScheduleScanConfig:
		// nolint:forcetypeassert
		scanConfig := scanConfigType.(*models.ByDaysScheduleScanConfig)
		interval, startTime = getIntervalAndStartTimeFromByDaysScheduleScanConfig(timeNow, scanConfig)
	case models.WeeklyScheduleScanConfig:
		// nolint:forcetypeassert
		scanConfig := scanConfigType.(*models.WeeklyScheduleScanConfig)
		interval, startTime = getIntervalAndStartTimeFromWeeklyScheduleScanConfig(timeNow, scanConfig)
	default:
		return nil, fmt.Errorf("unsupported schedule config type: %v", scanConfigType)
	}

	if interval <= 0 {
		return nil, fmt.Errorf("parameters validation failed. Interval=%v", interval)
	}

	return &Params{
		Interval:   interval,
		StartTime:  startTime,
		SingleScan: singleScan,
	}, nil
}

func (s *Scheduler) schedule(ctx context.Context, params *Params, scanConfig *models.ScanConfig) {

	startsAt := getStartsAt(time.Now().UTC(), params.StartTime, params.Interval)

	go s.spin(ctx, params, startsAt, scanConfig)
}

// get the next scan, that is after timeNow. if currentScanTime is already after timeNow, it will be return.
func getNextScanTime(timeNow, currentScanTime time.Time, interval time.Duration) time.Time {
	// if current scan time is before timeNow, jump to the next future scan time
	if currentScanTime.Before(timeNow) {
		// if scan time has passed in less than a second, start a scan now.
		timePassed := timeNow.Sub(currentScanTime)
		if timePassed < time.Second {
			return timeNow
		}
		remainingInterval := timePassed % interval
		if remainingInterval == 0 {
			currentScanTime = timeNow
		} else {
			currentScanTime = timeNow.Add(interval - remainingInterval)
		}
	}
	return currentScanTime
}

// get the time in Duration that the next scan should start at.
func getStartsAt(timeNow time.Time, startTime time.Time, interval time.Duration) time.Duration {
	nextScanTime := getNextScanTime(timeNow, startTime, interval)

	startsAt := nextScanTime.Sub(timeNow)

	return startsAt
}

func (s *Scheduler) spin(ctx context.Context, params *Params, startsAt time.Duration, scanConfig *models.ScanConfig) {
	log.Debugf("Starting a new scheduled scan. interval: %v, start time: %v, starts at: %v",
		params.Interval, params.StartTime, startsAt)
	singleScan := params.SingleScan
	interval := params.Interval

	timer := time.NewTimer(startsAt)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return
	case <-timer.C:
		go func() {
			scanDone := make(chan struct{})
			ticker := time.NewTicker(interval)
			defer ticker.Stop()
			for {
				if err := s.scan(ctx, scanConfig, scanDone); err != nil {
					log.Errorf("Failed to send scan: %v", err)
				}
				if singleScan {
					return
				}
				select {
				case <-ticker.C:
				case <-ctx.Done():
					log.Debugf("Received a stop signal...")
					return
				}
			}
		}()
	}
}

func (s *Scheduler) scan(ctx context.Context, scanConfig *models.ScanConfig, scanDone chan struct{}) error {
	// TODO: check if existing scan or a new scan
	targetInstances, scanID, err := s.initNewScan(ctx, scanConfig)
	if err != nil {
		return fmt.Errorf("failed to init new scan: %v", err)
	}

	scanner := _scanner.CreateScanner(s.scannerConfig, s.providerClient, s.backendClient, scanConfig, targetInstances, scanID)

	if err := scanner.Scan(ctx, scanDone); err != nil {
		return fmt.Errorf("failed to scan: %v", err)
	}

	return nil
}

// initNewScan Initialized a new scan, returns target instances and scan ID.
func (s *Scheduler) initNewScan(ctx context.Context, scanConfig *models.ScanConfig) ([]*types.TargetInstance, string, error) {
	instances, err := s.providerClient.Discover(ctx, scanConfig.Scope)
	if err != nil {
		return nil, "", fmt.Errorf("failed to discover instances to scan: %v", err)
	}

	targetInstances, err := s.createTargetInstances(ctx, instances)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get or create targets: %v", err)
	}

	now := time.Now().UTC()
	scan := &models.Scan{
		ScanConfigId:       scanConfig.Id,
		ScanFamiliesConfig: scanConfig.ScanFamiliesConfig,
		StartTime:          &now,
		TargetIDs:          getTargetIDs(targetInstances),
	}
	scanID, err := s.getOrCreateScan(ctx, scan)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get or create a scan: %v", err)
	}

	return targetInstances, scanID, nil
}

func getTargetIDs(targetInstances []*types.TargetInstance) *[]string {
	ret := make([]string, len(targetInstances))
	for i, targetInstance := range targetInstances {
		ret[i] = targetInstance.TargetID
	}

	return &ret
}

func (s *Scheduler) createTargetInstances(ctx context.Context, instances []types.Instance) ([]*types.TargetInstance, error) {
	targetInstances := make([]*types.TargetInstance, 0, len(instances))
	for i, instance := range instances {
		target, err := s.getOrCreateTarget(ctx, instance)
		if err != nil {
			return nil, fmt.Errorf("failed to get or create target. instanceID=%v: %v", instance.GetID(), err)
		}
		targetInstances = append(targetInstances, &types.TargetInstance{
			TargetID: *target.Id,
			Instance: instances[i],
		})
	}

	return targetInstances, nil
}

func (s *Scheduler) getOrCreateTarget(ctx context.Context, instance types.Instance) (*models.Target, error) {
	info := models.TargetType{}
	instanceProvider := models.AWS
	err := info.FromVMInfo(models.VMInfo{
		InstanceID:       utils.StringPtr(instance.GetID()),
		InstanceProvider: &instanceProvider,
		Location:         utils.StringPtr(instance.GetLocation()),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create VMInfo: %v", err)
	}
	resp, err := s.backendClient.PostTargetsWithResponse(ctx, models.Target{
		TargetInfo: &info,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to post target: %v", err)
	}
	switch resp.StatusCode() {
	case http.StatusCreated:
		if resp.JSON201 == nil {
			return nil, fmt.Errorf("failed to create a target: empty body")
		}
		return resp.JSON201, nil
	case http.StatusConflict:
		if resp.JSON409 == nil {
			return nil, fmt.Errorf("failed to create a target: empty body on conflict")
		}
		return resp.JSON409, nil
	default:
		if resp.JSONDefault != nil && resp.JSONDefault.Message != nil {
			return nil, fmt.Errorf("failed to post target. status code=%v: %v", resp.StatusCode(), resp.JSONDefault.Message)
		}
		return nil, fmt.Errorf("failed to post target. status code=%v", resp.StatusCode())
	}
}

// nolint:cyclop
func (s *Scheduler) getOrCreateScan(ctx context.Context, scan *models.Scan) (string, error) {
	resp, err := s.backendClient.PostScansWithResponse(ctx, *scan)
	if err != nil {
		return "", fmt.Errorf("failed to post a scan: %v", err)
	}
	switch resp.StatusCode() {
	case http.StatusCreated:
		if resp.JSON201 == nil {
			return "", fmt.Errorf("failed to create a scan: empty body")
		}
		if resp.JSON201.Id == nil {
			return "", fmt.Errorf("scan id is nil")
		}
		return *resp.JSON201.Id, nil
	case http.StatusConflict:
		if resp.JSON409 == nil {
			return "", fmt.Errorf("failed to create a scan: empty body on conflict")
		}
		if resp.JSON409.Id == nil {
			return "", fmt.Errorf("scan id on conflict is nil")
		}
		return *resp.JSON409.Id, nil
	default:
		if resp.JSONDefault != nil && resp.JSONDefault.Message != nil {
			return "", fmt.Errorf("failed to post scan. status code=%v: %v", resp.StatusCode(), resp.JSONDefault.Message)
		}
		return "", fmt.Errorf("failed to post scan. status code=%v", resp.StatusCode())
	}
}
