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

package gcp

import (
	"cloud.google.com/go/compute/apiv1/computepb"
	"context"
	"fmt"
	"github.com/openclarity/vmclarity/testenv/utils"
	dockerhelper "github.com/openclarity/vmclarity/testenv/utils/docker"
)

func (e *GCPEnv) afterSetUp(ctx context.Context) error {
	req := &computepb.GetInstanceRequest{
		Project:  ProjectID,
		Zone:     Zone,
		Instance: "vmclarity-" + e.envName + "-server",
	}
	server, err := e.instancesClient.Get(ctx, req)
	e.serverIP = server.NetworkInterfaces[0].AccessConfigs[0].NatIP

	e.sshPortForwardInput = &utils.SSHForwardInput{
		PrivateKey:    e.sshKeyPair.PrivateKey,
		User:          DefaultRemoteUser,
		Host:          *e.serverIP,
		Port:          utils.DefaultSSHPort,
		LocalPort:     8080, //nolint:gomnd
		RemoteAddress: "localhost",
		RemotePort:    80, //nolint:gomnd
	}

	e.sshPortForward, err = utils.NewSSHPortForward(e.sshPortForwardInput)
	if err != nil {
		return fmt.Errorf("failed to setup SSH port forwarding: %w", err)
	}

	if err = e.sshPortForward.Start(context.Background()); err != nil { //nolint:contextcheck
		return fmt.Errorf("failed to wait for SSH port to become ready: %w", err)
	}

	clientOpts, err := dockerhelper.ClientOptsWithSSHConn(ctx, e.workDir, e.sshKeyPair, e.sshPortForwardInput)
	if err != nil {
		return fmt.Errorf("failed to get options for docker client: %w", err)
	}

	e.DockerHelper, err = dockerhelper.New(clientOpts)
	if err != nil {
		return fmt.Errorf("failed to create Docker helper: %w", err)
	}

	err = e.DockerHelper.WaitForDockerReady(ctx)
	if err != nil {
		return fmt.Errorf("failed to check if Docker client is ready: %w", err)
	}

	return nil
}
