// Copyright © 2024 Cisco Systems, Inc. and its affiliates.
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

package binary

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/openclarity/vmclarity/plugins/runner/internal/runtimehandler"
	"github.com/openclarity/vmclarity/plugins/runner/types"
	"github.com/openclarity/vmclarity/plugins/sdk-go/plugin"
	"github.com/openclarity/vmclarity/scanner/utils/containerrootfs"
)

type binaryRuntimeHandler struct {
	config types.PluginConfig

	cmd       *exec.Cmd
	cmdStdOut bytes.Buffer
	cmdStdErr bytes.Buffer

	pluginDir            string
	inputDirMountPoint   string
	pluginServerEndpoint string
	ready                bool
}

func New(ctx context.Context, config types.PluginConfig) (runtimehandler.PluginRuntimeHandler, error) {
	return &binaryRuntimeHandler{
		config: config,
	}, nil
}

func (h *binaryRuntimeHandler) Start(ctx context.Context) error {
	h.pluginDir = filepath.Join(os.TempDir(), "vmclarity-plugins", h.config.Name)

	image, cleanup, err := containerrootfs.GetImageWithCleanup(ctx, h.config.ImageName)
	if err != nil {
		return fmt.Errorf("unable to get image(%s): %w", h.config.ImageName, err)
	}
	defer cleanup()

	if _, err := os.Stat(h.pluginDir); os.IsNotExist(err) {
		err = containerrootfs.ToDirectory(ctx, image, h.pluginDir)
		if err != nil {
			return fmt.Errorf("unable to extract image(%s): %w", h.config.ImageName, err)
		}
	}

	// Mount input from host
	// /tmp/vmclarity-plugins/kics + /input + /host-dir-to-scan
	h.inputDirMountPoint = filepath.Join(h.pluginDir, runtimehandler.RemoteScanInputDirOverride, h.config.InputDir)
	err = os.MkdirAll(h.inputDirMountPoint, 0755)
	if err != nil {
		return fmt.Errorf("unable to create directory for mount point: %w", err)
	}

	// nolint:typecheck
	err = syscall.Mount(h.config.InputDir, h.inputDirMountPoint, "", syscall.MS_BIND, "")
	if err != nil {
		return fmt.Errorf("unable to mount input directory (%s - %s): %w", h.config.InputDir, h.inputDirMountPoint, err)
	}

	// https://lwn.net/Articles/281157/
	// "the read-only attribute can only be added with a remount operation afterward"
	err = syscall.Mount(h.config.InputDir, h.inputDirMountPoint, "", syscall.MS_BIND|syscall.MS_REMOUNT|syscall.MS_RDONLY, "")
	if err != nil {
		return fmt.Errorf("unable to remount input directory as read-only (%s - %s): %w", h.config.InputDir, h.inputDirMountPoint, err)
	}

	// TODO: handle panics - unmount

	// Determine entrypoint or command to execute
	var args []string
	if len(image.Metadata.Config.Config.Entrypoint) > 0 {
		args = append(image.Metadata.Config.Config.Entrypoint[0:], image.Metadata.Config.Config.Cmd...)
	} else if len(image.Metadata.Config.Config.Cmd) > 0 {
		args = image.Metadata.Config.Config.Cmd[0:]
	} else {
		return fmt.Errorf("no entrypoint or command found in the config")
	}

	// Find a port
	openPortListener, err := net.Listen("tcp", ":0")
	if err != nil {
		return fmt.Errorf("unable to find port")
	}
	port := openPortListener.Addr().(*net.TCPAddr).Port

	h.pluginServerEndpoint = fmt.Sprintf("http://127.0.0.1:%d", port)

	// Set environment variables
	env := image.Metadata.Config.Config.Env
	env = append(env, fmt.Sprintf("%s=%s", plugin.EnvListenAddress, h.pluginServerEndpoint))

	// Set workdir
	workDir := image.Metadata.Config.Config.WorkingDir
	if workDir == "" {
		workDir = "/"
	}

	// Initialize command
	h.cmd = exec.CommandContext(ctx, args[0], args[1:]...)
	h.cmd.Env = env
	h.cmd.Dir = workDir

	h.cmd.SysProcAttr = &syscall.SysProcAttr{
		Chroot: h.pluginDir,
	}

	h.cmd.Stdin = &bytes.Buffer{}
	h.cmd.Stdout = &h.cmdStdOut
	h.cmd.Stderr = &h.cmdStdErr

	// Start command
	openPortListener.Close()
	err = h.cmd.Start()
	if err != nil {
		return fmt.Errorf("unable to start process: %w", err)
	}

	// Waiting for command to be finished in the background
	go func() {
		h.cmd.Wait()
	}()

	h.ready = true

	return nil
}

func (h *binaryRuntimeHandler) Ready() (bool, error) {
	if h.cmd == nil {
		return false, fmt.Errorf("plugin process is not running")
	}

	return h.ready, nil
}

func (h *binaryRuntimeHandler) GetPluginServerEndpoint(ctx context.Context) (string, error) {
	return h.pluginServerEndpoint, nil
}

func (h *binaryRuntimeHandler) Logs(ctx context.Context) (io.ReadCloser, error) {
	if h.cmd == nil {
		return nil, fmt.Errorf("plugin process is not running")
	}

	reader := io.MultiReader(&h.cmdStdOut, &h.cmdStdErr)
	return io.NopCloser(reader), nil
}

func (h *binaryRuntimeHandler) Result(ctx context.Context) (io.ReadCloser, error) {
	f, err := os.Open(filepath.Join(h.pluginDir, runtimehandler.RemoteScanResultFileOverride))
	if err != nil {
		return nil, fmt.Errorf("unable to open result file: %w", err)
	}

	return f, nil
}

func (h *binaryRuntimeHandler) Remove(ctx context.Context) error {
	var removeErr error

	if h.cmd.ProcessState != nil {
		if h.cmd.ProcessState.Exited() == false {
			if err := h.cmd.Process.Kill(); err != nil {
				removeErr = multierror.Append(removeErr, fmt.Errorf("failed to kill plugin process: %w", err))
			}
		}
	}

	// Unmount input directory
	if err := syscall.Unmount(h.inputDirMountPoint, 0); err != nil {
		removeErr = multierror.Append(removeErr, fmt.Errorf("failed to kill plugin process: %w", err))
	}

	return removeErr
}
