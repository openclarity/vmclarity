// Copyright © 2022 Cisco Systems, Inc. and its affiliates.
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

package creds

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/containers/image/v5/docker/reference"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/authn/k8schain"
	"github.com/google/go-containerregistry/pkg/name"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

func GetAuthConfig(ipsPath string, imageName string) (*authn.AuthConfig, error) {
	named, err := reference.ParseNormalizedNamed(imageName)
	if err != nil {
		return nil, fmt.Errorf("failed to normalized image name: %w", err)
	}

	keychain, err := newKeyChain(context.TODO(), ipsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create keychain: %w", err)
	}

	repository, err := name.NewRepository(reference.FamiliarName(named))
	if err != nil {
		return nil, fmt.Errorf("failed to create repository: %w", err)
	}

	authenticator, err := keychain.Resolve(repository)
	if err != nil {
		return nil, fmt.Errorf("failed to create authenticator: %w", err)
	}

	authorization, err := authenticator.Authorization()
	if err != nil {
		return nil, fmt.Errorf("failed to create authorization: %w", err)
	}

	return authorization, nil
}

func newKeyChain(ctx context.Context, ipsPath string) (authn.Keychain, error) {
	if ipsPath == "" {
		keychain, err := k8schain.NewNoClient(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to create no client keychain: %w", err)
		}

		return keychain, nil
	}

	secrets, err := readImagePullSecrets(ipsPath)
	if err != nil {
		return nil, fmt.Errorf("fail to read image pull secrets: %w", err)
	}

	keychain, err := k8schain.NewFromPullSecrets(ctx, secrets)
	if err != nil {
		return nil, fmt.Errorf("unable to load keychain from image pull secrets: %w", err)
	}

	return keychain, nil
}

func readImagePullSecrets(ipsPath string) ([]corev1.Secret, error) {
	secrets := []corev1.Secret{}
	files, err := os.ReadDir(ipsPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read path %s: %w", ipsPath, err)
	}

	for _, file := range files {
		// We expect directories for each secret in ipsPath
		if !file.IsDir() {
			continue
		}

		secretPath := path.Join(ipsPath, file.Name())
		secretType, secretDataKey, err := determineSecretTypeAndKey(secretPath)
		if err != nil {
			return nil, fmt.Errorf("unable to determine type of secret %s: %w", file.Name(), err)
		}

		if secretType != corev1.SecretTypeDockerConfigJson && secretType != corev1.SecretTypeDockercfg {
			log.Warnf("Secret %s is not a supported image pull secret type, ignoring.", file.Name())
			continue
		}

		secretFilePath := path.Join(ipsPath, file.Name(), secretDataKey)
		secretDataBytes, err := os.ReadFile(secretFilePath)
		if err != nil {
			return nil, fmt.Errorf("unable to read secret file %s: %w", secretFilePath, err)
		}

		secrets = append(secrets, corev1.Secret{
			Type: secretType,
			Data: map[string][]byte{
				secretDataKey: secretDataBytes,
			},
		})
	}

	return secrets, nil
}

func determineSecretTypeAndKey(secretPath string) (corev1.SecretType, string, error) {
	var unsetSecretType corev1.SecretType

	secretFiles, err := os.ReadDir(secretPath)
	if err != nil {
		return unsetSecretType, "", fmt.Errorf("unable to read secret directory %s: %w", secretPath, err)
	}

	for _, secretFile := range secretFiles {
		// We only want files at this point
		if secretFile.IsDir() {
			continue
		}

		switch secretFile.Name() {
		case corev1.DockerConfigJsonKey:
			return corev1.SecretTypeDockerConfigJson, corev1.DockerConfigJsonKey, nil
		case corev1.DockerConfigKey:
			return corev1.SecretTypeDockercfg, corev1.DockerConfigKey, nil
		}
	}

	return unsetSecretType, "", nil
}
