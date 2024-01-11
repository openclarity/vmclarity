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

package image_helper // nolint:revive,stylecheck

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/anchore/stereoscope/pkg/image"
	"github.com/containers/image/v5/docker/reference"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	containerregistry_v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/daemon"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	log "github.com/sirupsen/logrus"

	"github.com/openclarity/vmclarity/pkg/shared/config"
)

// FsLayerCommand represents a history command of a layer in a docker image.
type FsLayerCommand struct {
	Command string
	Layer   string
}

func GetHashFromRepoDigest(repoDigests []string, imageName string) string {
	if len(repoDigests) == 0 {
		return ""
	}

	normalizedName, err := reference.ParseNormalizedNamed(imageName)
	if err != nil {
		log.Errorf("Failed to parse image name %s to normalized named: %v", imageName, err)
		return ""
	}
	familiarName := reference.FamiliarName(normalizedName)
	// iterating over RepoDigests and use RepoDigest which match to imageName
	for _, repoDigest := range repoDigests {
		normalizedRepoDigest, err := reference.ParseNormalizedNamed(repoDigest)
		if err != nil {
			log.Errorf("Failed to parse repoDigest %s, %v", repoDigest, err)
			return ""
		}
		// RepoDigests can be different based on the registry
		//        ],
		//        "RepoDigests": [
		//            "debian@sha256:2906804d2a64e8a13a434a1a127fe3f6a28bf7cf3696be4223b06276f32f1f2d",
		//            "poke/debian@sha256:a4c378901a2ba14fd331e96a49101556e91ed592d5fd68ba7405fdbf9b969e61",
		//            "poke/testdebian@sha256:a4c378901a2ba14fd331e96a49101556e91ed592d5fd68ba7405fdbf9b969e61"
		//        ],
		// Check which RegoDigest should be used
		if reference.FamiliarName(normalizedRepoDigest) == familiarName {
			return normalizedRepoDigest.(reference.Digested).Digest().Encoded() // nolint:forcetypeassert
		}
	}
	return ""
}

// fetchFsCommands retrieves information about image layers commands.
func fetchFsCommands(img containerregistry_v1.Image) ([]*FsLayerCommand, error) {
	configFile, err := img.RawConfigFile()
	if err != nil {
		return nil, fmt.Errorf("failed to get raw config file: %w", err)
	}

	var conf containerregistry_v1.ConfigFile
	if err = json.Unmarshal(configFile, &conf); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	if log.IsLevelEnabled(log.DebugLevel) {
		confB, err := json.Marshal(conf)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal config: %w", err)
		}
		log.Debugf("Image config: %s", confB)
	}

	commands := getCommands(&conf)

	layers, err := img.Layers()
	if err != nil {
		return nil, fmt.Errorf("failed to get layers: %w", err)
	}

	if len(layers) != len(commands) {
		log.Infof("Number of fs layers (%v) doesn't match the number of fs history entries (%v) - setting empty commands", len(layers), len(commands))
		commands = make([]string, len(layers))
	}

	fsLayerCommands, err := createFsLayerCommands(layers, commands)
	if err != nil {
		return nil, fmt.Errorf("failed to create fs layer commands: %w", err)
	}

	if log.IsLevelEnabled(log.DebugLevel) {
		fsLayerCommandsB, err := json.Marshal(fsLayerCommands)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal layer commands: %w", err)
		}
		log.Debugf("Layer commands: %s", fsLayerCommandsB)
	}

	return fsLayerCommands, nil
}

func createFsLayerCommands(layers []containerregistry_v1.Layer, commands []string) ([]*FsLayerCommand, error) {
	layerCommands := make([]*FsLayerCommand, len(layers))

	for i, layer := range layers {
		layerDiffID, err := layer.DiffID() // specifies the Hash of the uncompressed layer
		if err != nil {
			return nil, fmt.Errorf("failed to get layer diffID: %w", err)
		}
		layerCommands[i] = &FsLayerCommand{
			Command: commands[i],
			Layer:   layerDiffID.Hex,
		}
	}

	return layerCommands, nil
}

func getCommands(conf *containerregistry_v1.ConfigFile) []string {
	// nolint:prealloc
	var commands []string
	for i, layerHistory := range conf.History {
		if layerHistory.EmptyLayer {
			log.Infof("Skipping empty layer (%v): %+v", i, layerHistory)
			continue
		}
		commands = append(commands, stripDockerMetaFromCommand(layerHistory.CreatedBy))
	}
	return commands
}

// Strips Dockerfile generation info from layer commands. e.g: "/bin/sh -c #(nop) CMD [/bin/bash]" -> "CMD [/bin/bash]".
func stripDockerMetaFromCommand(command string) string {
	ret := strings.TrimSpace(strings.TrimPrefix(command, "/bin/sh -c #(nop)"))
	ret = strings.TrimSpace(strings.TrimPrefix(ret, "/bin/sh -c"))
	return ret
}

func getV1Image(imageName string, registryOptions *image.RegistryOptions, localImage bool) (containerregistry_v1.Image, error) {
	ref, err := name.ParseReference(imageName, prepareReferenceOptions(registryOptions)...)
	if err != nil {
		return nil, fmt.Errorf("unable to parse registry reference=%q: %w", imageName, err)
	}

	switch localImage {
	case true:
		img, err := daemon.Image(ref, daemon.WithUnbufferedOpener())
		if err != nil {
			return nil, fmt.Errorf("failed to get image from daemon: %w", err)
		}
		return img, nil
	default:
		log.Debugf("pulling image info directly from registry image=%q", imageName)
		img, err := remote.Image(ref, prepareRemoteOptions(ref, registryOptions)...)
		if err != nil {
			return nil, fmt.Errorf("failed to get image from registry: %w", err)
		}
		return img, nil
	}
}

func prepareReferenceOptions(registryOptions *image.RegistryOptions) []name.Option {
	var options []name.Option
	if registryOptions != nil && registryOptions.InsecureUseHTTP {
		options = append(options, name.Insecure)
	}
	return options
}

func prepareRemoteOptions(ref name.Reference, registryOptions *image.RegistryOptions) []remote.Option {
	opts := make([]remote.Option, 0)
	if registryOptions == nil {
		// use the Keychain specified from a docker config file.
		log.Debugf("no registry credentials configured, using the default keychain")
		opts = append(opts, remote.WithAuthFromKeychain(authn.DefaultKeychain))
		return opts
	}

	if registryOptions.InsecureSkipTLSVerify {
		t := &http.Transport{
			// nolint: gosec
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		opts = append(opts, remote.WithTransport(t))
	}

	// note: the authn.Authenticator and authn.Keychain options are mutually exclusive, only one may be provided.
	// If no explicit authenticator can be found, then fallback to the keychain.
	authenticator := registryOptions.Authenticator(ref.Context().RegistryStr())
	if authenticator != nil {
		opts = append(opts, remote.WithAuth(authenticator))
	}

	return opts
}

func GetImageLayerCommands(imageName string, sharedConf *config.Config) ([]*FsLayerCommand, error) {
	registryOptions := config.CreateRegistryOptions(sharedConf.Registry)
	img, err := getV1Image(imageName, registryOptions, sharedConf.LocalImageScan)
	if err != nil {
		return nil, fmt.Errorf("failed to get v1.image=%s: %w", imageName, err)
	}
	layerCommands, err := fetchFsCommands(img)
	if err != nil {
		return nil, fmt.Errorf("failed to get layer commands from image=%s: %w", imageName, err)
	}
	return layerCommands, nil
}

func GetHashFromRepoDigestsOrImageID(repoDigests []string, imageID string, imageName string) (string, error) {
	if imageID == "" && len(repoDigests) == 0 {
		return "", fmt.Errorf("RepoDigest and ImageID are missing")
	}

	hash := GetHashFromRepoDigest(repoDigests, imageName)
	if hash == "" {
		// set hash using ImageID (https://github.com/opencontainers/image-spec/blob/main/config.md#imageid) if repo digests are missing
		// image ID is represented as a hexadecimal encoding of 256 bits, e.g., sha256:a9561eb1b190625c9adb5a9513e72c4dedafc1cb2d4c5236c9a6957ec7dfd5a9
		// we need only the hash
		_, h, found := strings.Cut(imageID, ":")
		if found {
			hash = h
		} else {
			hash = imageID
		}
	}
	return hash, nil
}
