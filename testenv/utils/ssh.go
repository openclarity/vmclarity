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

package utils

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
)

type SSHKeyPair struct {
	PublicKeyFile  string
	PrivateKeyFile string
	Temporary      bool
}

// GenerateSSHKeyPair generates a new SSH key pair.
func GenerateSSHKeyPair() (*SSHKeyPair, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}

	privateKeyFile, err := os.Create("id_rsa_testenv")
	defer privateKeyFile.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to create private key file: %w", err)
	}

	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	if err := pem.Encode(privateKeyFile, privateKeyPEM); err != nil {
		return nil, fmt.Errorf("failed to encode private key: %w", err)
	}

	publicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create public key: %w", err)
	}

	publicKeyFile := "id_rsa_testenv.pub"
	err = os.WriteFile(publicKeyFile, ssh.MarshalAuthorizedKey(publicKey), 0655)
	if err != nil {
		return nil, fmt.Errorf("failed to write public key file: %w", err)
	}

	return &SSHKeyPair{
		PublicKeyFile:  publicKeyFile,
		PrivateKeyFile: privateKeyFile.Name(),
		Temporary:      true,
	}, nil
}

// Run SSH tunnel to remote VMClarity server.
func RunSSHTunnel(ctx context.Context, privateKeyFile, remoteHost, remotePort, localPort string) {
	logger := GetLoggerFromContextOrDiscard(ctx)

	// Read the private key file.
	key, err := os.ReadFile(privateKeyFile)
	if err != nil {
		logger.Fatalf("failed to read private key: %s\n", err)
		return
	}

	// Create Signer from private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		logger.Fatalf("failed to parse private key: %s\n", err)
		return
	}

	// Create SSH client config.
	config := ssh.ClientConfig{
		User: "ubuntu",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		// TODO(paralta) Review this insecure host key callback, which accepts any host key.
		HostKeyCallback:   ssh.InsecureIgnoreHostKey(), //nolint:gosec
		HostKeyAlgorithms: []string{ssh.KeyAlgoED25519},
	}

	// Add port to remote host.
	remoteAddress := fmt.Sprintf("%s:%s", remoteHost, "22")

	// Dial the remote server.
	client, err := ssh.Dial("tcp", remoteAddress, &config)
	if err != nil {
		logger.Fatalf("failed to dial: %s\n", err)
	}
	defer client.Close()

	// Listen on local port.
	listener, err := net.Listen("tcp", "localhost:"+localPort)
	if err != nil {
		logger.Fatalf("failed to listen: %s\n", err)
	}
	defer listener.Close()

	for {
		// Accept local connection.
		local, err := listener.Accept()
		if err != nil {
			logger.Fatalf("failed to accept: %s\n", err)
		}

		// Dial remote server.
		remote, err := client.Dial("tcp", "localhost:"+remotePort)
		if err != nil {
			logger.Fatalf("failed to dial: %s\n", err)
		}

		// Run tunnel between local and remote connections.
		runTunnel(ctx, local, remote)
	}
}

// runTunnel runs a tunnel between two connections.
func runTunnel(ctx context.Context, local, remote net.Conn) {
	logger := GetLoggerFromContextOrDiscard(ctx)

	defer local.Close()
	defer remote.Close()
	done := make(chan struct{}, 2) //nolint:gomnd

	go func() {
		_, err := io.Copy(local, remote)
		if err != nil {
			logger.Fatalf("failed to copy data from remote to local: %s\n", err)
		}
		done <- struct{}{}
	}()

	go func() {
		_, err := io.Copy(remote, local)
		if err != nil {
			logger.Fatalf("failed to copy data from local to remote: %s\n", err)
		}
		done <- struct{}{}
	}()

	<-done
}
