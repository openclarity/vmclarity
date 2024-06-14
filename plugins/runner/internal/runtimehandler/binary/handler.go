package binary

import (
	"context"
	"fmt"
	"github.com/openclarity/vmclarity/utils/command"
	"io"
	"os"
	"path/filepath"
	"syscall"

	"github.com/anchore/stereoscope"
	"github.com/openclarity/vmclarity/plugins/runner/internal/runtimehandler"
	"github.com/openclarity/vmclarity/plugins/runner/types"
	"github.com/openclarity/vmclarity/scanner/utils/containerrootfs"
)

type binaryRuntimeHandler struct {
	config types.PluginConfig

	commandState *command.State
}

func New(ctx context.Context, config types.PluginConfig) (runtimehandler.PluginRuntimeHandler, error) {
	return &binaryRuntimeHandler{
		config: config,
	}, nil
}

func (h *binaryRuntimeHandler) Start(ctx context.Context) error {
	pluginDir := filepath.Join(os.TempDir(), "vmclarity-plugins", h.config.Name)

	image, err := stereoscope.GetImage(ctx, h.config.ImageName)
	if err != nil {
		return fmt.Errorf("unable to get image(%s): %w", h.config.ImageName, err)
	}

	// TODO: check if directory exists
	err = containerrootfs.ToDirectory(ctx, h.config.ImageName, pluginDir)
	if err != nil {
		return fmt.Errorf("unable to extract image(%s): %w", h.config.ImageName, err)
	}

	// or bind mount?
	hostFsDir := filepath.Join(pluginDir, "hostfs")
	err = os.Symlink("/", hostFsDir)
	if err != nil {
		return fmt.Errorf("unable to create hostfs symlink(%s): %w", h.config.ImageName, err)
	}

	//os.RemoveAll(filepath.Join(pluginDir, "dev"))
	//err = os.Symlink("/dev", filepath.Join(pluginDir, "dev"))
	//if err != nil {
	//	return fmt.Errorf("unable to create hostfs symlink(%s): %w", h.config.ImageName, err)
	//}

	err = syscall.Chroot(pluginDir)
	if err != nil {
		return fmt.Errorf("unable to call chroot: %w", err)
	}

	var args []string
	if len(image.Metadata.Config.Config.Entrypoint) > 0 {
		args = append(image.Metadata.Config.Config.Entrypoint[0:], image.Metadata.Config.Config.Cmd...)
	} else if len(image.Metadata.Config.Config.Cmd) > 0 {
		args = image.Metadata.Config.Config.Cmd[0:]
	} else {
		return fmt.Errorf("no entrypoint or command found in the config")
	}

	// lock os thread?
	env := image.Metadata.Config.Config.Env
	//workDir := image.Metadata.Config.Config.WorkingDir
	err = syscall.Exec(args[0], args, env)
	if err != nil {
		return fmt.Errorf("unable to start process: %w", err)
	}

	/*cmd := command.Command{
		Cmd:     arg0,
		Args:    args,
		Env:     env,
		WorkDir: workDir,
	}

	h.commandState, err = cmd.Start(ctx)
	if err != nil {
		return fmt.Errorf("unable to start process: %w", err)
	}

	err = h.commandState.Wait()
	if err != nil {
		return fmt.Errorf("unable to wait process: %w", err)
	}*/

	return nil
}

func (h *binaryRuntimeHandler) Ready() (bool, error) {
	panic("not implemented")

	// check plugin process state

	return false, nil
}

func (h *binaryRuntimeHandler) GetPluginServerEndpoint(ctx context.Context) (string, error) {
	panic("not implemented")

	// find open port to listen

	return "", nil
}

func (h *binaryRuntimeHandler) Logs(ctx context.Context) (io.ReadCloser, error) {
	panic("not implemented")

	return nil, nil
}

func (h *binaryRuntimeHandler) Result(ctx context.Context) (io.ReadCloser, error) {
	panic("not implemented")

	return nil, nil
}

func (h *binaryRuntimeHandler) Remove(ctx context.Context) error {
	panic("not implemented")

	return nil
}
