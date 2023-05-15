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

//nolint:wrapcheck
package cli

import (
	"context"
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"

	"github.com/openclarity/vmclarity/cli/pkg/mount"
	"github.com/openclarity/vmclarity/cli/pkg/presenter"
	"github.com/openclarity/vmclarity/cli/pkg/state"
	"github.com/openclarity/vmclarity/shared/pkg/families"
	"github.com/openclarity/vmclarity/shared/pkg/families/types"
	"github.com/openclarity/vmclarity/shared/pkg/log"
)

const (
	fsTypeExt4 = "ext4"
	fsTypeXFS  = "xfs"
)

type CLI struct {
	state.Manager
	presenter.Presenter

	FamiliesConfig *families.Config
}

func (c *CLI) FamilyStarted(ctx context.Context, famType types.FamilyType) error {
	return c.Manager.MarkFamilyScanInProgress(ctx, famType)
}

func (c *CLI) FamilyFinished(ctx context.Context, res families.FamilyResult) error {
	return c.Presenter.ExportFamilyResult(ctx, res)
}

func (c *CLI) MountVolumes(ctx context.Context) ([]string, error) {
	var mountPoints []string

	logger := log.GetLoggerFromContextOrDiscard(ctx)

	devices, err := mount.ListBlockDevices()
	if err != nil {
		return nil, fmt.Errorf("failed to list block devices: %v", err)
	}
	for _, device := range devices {
		// if the device is not mounted and of a supported filesystem type,
		// we assume it belongs to the attached volume, so we mount it.
		if device.MountPoint == "" && isSupportedFS(device.FilesystemType) {
			mountDir := "/mnt/snapshot" + uuid.NewV4().String()

			if err := device.Mount(mountDir); err != nil {
				return nil, fmt.Errorf("failed to mount device: %v", err)
			}
			logger.Infof("Device %v on %v is mounted", device.DeviceName, mountDir)
			mountPoints = append(mountPoints, mountDir)
		}
		if ctx.Err() != nil {
			return mountPoints, fmt.Errorf("failed to mount block devices: %w", ctx.Err())
		}
	}
	return mountPoints, nil
}

// WatchForAbort is responsible for watching for abort events triggered and invoking the provided cancel function to mark
// the ctx context cancelled. Abort events are fired by:
// * ScanResult state is changed to ABORTED
// * context deadline defined by timeout parameter is reached
func (c *CLI) WatchForAbort(ctx context.Context, cancel context.CancelFunc, interval time.Duration, timeout time.Duration) {
	go func() {
		ctx, timeoutCancel := context.WithTimeout(ctx, timeout)
		timer := time.NewTicker(interval)

		defer func() {
			timer.Stop()
			timeoutCancel()
			cancel()
		}()

		logger := log.GetLoggerFromContextOrDiscard(ctx)

		for {
			select {
			case <-timer.C:
				aborted, err := c.IsAborted(ctx)
				if err != nil {
					logger.Errorf("Failed to retrieve scan result state: %v", err)
				}
				if aborted {
					return
				}
			case <-ctx.Done():
				logger.Debugf("Stop watching for abort event as context is cancelled")
				return
			}
		}
	}()
}

func isSupportedFS(fs string) bool {
	switch fs {
	case fsTypeExt4, fsTypeXFS:
		return true
	}
	return false
}
