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

package sshtopology

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/openclarity/kubeclarity/shared/pkg/job_manager"
	"github.com/openclarity/kubeclarity/shared/pkg/utils"
	log "github.com/sirupsen/logrus"

	"github.com/openclarity/vmclarity/pkg/shared/families/infofinder/types"
	sharedUtils "github.com/openclarity/vmclarity/pkg/shared/utils"
)

const ScannerName = "sshTopology"

type Scanner struct {
	name       string
	logger     *log.Entry
	config     types.SSHTopologyConfig
	resultChan chan job_manager.Result
}

func New(c job_manager.IsConfig, logger *log.Entry, resultChan chan job_manager.Result) job_manager.Job {
	conf := c.(types.ScannersConfig) // nolint:forcetypeassert
	return &Scanner{
		name:       ScannerName,
		logger:     logger.Dup().WithField("scanner", ScannerName),
		config:     conf.SSHTopology,
		resultChan: resultChan,
	}
}

func (s *Scanner) Run(sourceType utils.SourceType, userInput string) error {
	go func() {
		s.logger.Debugf("Running with input=%v and source type=%v", userInput, sourceType)
		retResults := types.ScannerResult{
			ScannerName: ScannerName,
		}

		// Validate this is an input type supported by the scanner,
		// otherwise return skipped.
		if !s.isValidInputType(sourceType) {
			s.sendResults(retResults, nil)
			return
		}

		var retErr error
		homeUserDirs, err := getHomeUserDirs(userInput)
		if err != nil {
			// Collect the error and continue.
			retErr = errors.Join(retErr, fmt.Errorf("failed to get home user dirs: %v", err))
		}
		s.logger.Debugf("Found home user dirs %+v", homeUserDirs)

		// The jobs are:
		// getSSHDaemonKeysFingerprints
		// getSSHPrivateKeysFingerprints (for each user folder)
		// getSSHAuthorizedKeysFingerprints (for each user folder)
		// getSSHKnownHostsFingerprints (for each user folder)
		jobsCount := 1 + 3*len(homeUserDirs)
		errs := make(chan error, jobsCount)
		fingerprintsChan := make(chan []types.Info, jobsCount)

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			if sshDaemonKeysFingerprints, err := s.getSSHDaemonKeysFingerprints(userInput); err != nil {
				errs <- fmt.Errorf("failed to get ssh daemon keys: %v", err)
			} else {
				fingerprintsChan <- sshDaemonKeysFingerprints
			}
		}()

		for i := range homeUserDirs {
			dir := homeUserDirs[i]

			wg.Add(1)
			go func() {
				defer wg.Done()
				if sshPrivateKeysFingerprints, err := s.getSSHPrivateKeysFingerprints(dir); err != nil {
					errs <- fmt.Errorf("failed to get ssh private keys: %v", err)
				} else {
					fingerprintsChan <- sshPrivateKeysFingerprints
				}
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				if sshAuthorizedKeysFingerprints, err := s.getSSHAuthorizedKeysFingerprints(dir); err != nil {
					errs <- fmt.Errorf("failed to get ssh authorized keys: %v", err)
				} else {
					fingerprintsChan <- sshAuthorizedKeysFingerprints
				}
			}()

			wg.Add(1)
			go func() {
				defer wg.Done()
				if sshKnownHostsFingerprints, err := s.getSSHKnownHostsFingerprints(dir); err != nil {
					errs <- fmt.Errorf("failed to get ssh known hosts: %v", err)
				} else {
					fingerprintsChan <- sshKnownHostsFingerprints
				}
			}()
		}

		wg.Wait()
		close(errs)
		close(fingerprintsChan)

		for e := range errs {
			if e != nil {
				retErr = errors.Join(retErr, e)
			}
		}
		if retErr != nil {
			retResults.Error = retErr
		}

		for fingerprints := range fingerprintsChan {
			retResults.Infos = append(retResults.Infos, fingerprints...)
		}

		if len(retResults.Infos) > 0 && retResults.Error != nil {
			// Since we have findings we want to share what we've got and only prints the errors here.
			// Maybe we need to support to send both errors and findings in a higher level.
			s.logger.Error(retResults.Error)
			retResults.Error = nil
		}

		s.sendResults(retResults, nil)
	}()

	return nil
}

func getHomeUserDirs(rootDir string) ([]string, error) {
	var dirs []string
	homeDirPath := path.Join(rootDir, "home")
	files, err := os.ReadDir(homeDirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read dir (%v): %v", homeDirPath, err)
	}

	for _, f := range files {
		if f.IsDir() {
			dirs = append(dirs, path.Join(homeDirPath, f.Name()))
		}
	}

	return dirs, nil
}

func (s *Scanner) getSSHDaemonKeysFingerprints(rootPath string) ([]types.Info, error) {
	paths, err := s.getPrivateKeysPaths(path.Join(rootPath, "/etc/ssh"), false)
	if err != nil {
		return nil, err
	}
	s.logger.Debugf("Found ssh daemon private keys paths %+v", paths)

	fingerprints, err := s.getFingerprints(paths, types.SSHDaemonKeys)
	if err != nil {
		return nil, err
	}
	s.logger.Debugf("Found ssh daemon private keys fingerprints %+v", fingerprints)

	return fingerprints, nil
}

func (s *Scanner) getSSHPrivateKeysFingerprints(homeUserDir string) ([]types.Info, error) {
	paths, err := s.getPrivateKeysPaths(homeUserDir, true)
	if err != nil {
		return nil, err
	}
	s.logger.Debugf("Found ssh private keys paths %+v", paths)

	infos, err := s.getFingerprints(paths, types.SSHPrivateKeys)
	if err != nil {
		return nil, err
	}
	s.logger.Debugf("Found ssh private keys fingerprints %+v", infos)

	return infos, nil
}

func (s *Scanner) getSSHAuthorizedKeysFingerprints(homeUserDir string) ([]types.Info, error) {
	infos, err := s.getFingerprints([]string{path.Join(homeUserDir, ".ssh/authorized_keys")}, types.SSHAuthorizedKeys)
	if err != nil {
		return nil, err
	}
	s.logger.Debugf("Found ssh authorized keys fingerprints %+v", infos)

	return infos, nil
}

func (s *Scanner) getSSHKnownHostsFingerprints(homeUserDir string) ([]types.Info, error) {
	infos, err := s.getFingerprints([]string{path.Join(homeUserDir, ".ssh/known_hosts")}, types.SSHKnownHosts)
	if err != nil {
		return nil, err
	}
	s.logger.Debugf("Found ssh known hosts fingerprints %+v", infos)

	return infos, nil
}

func (s *Scanner) getFingerprints(paths []string, infoType types.InfoType) ([]types.Info, error) {
	var infos []types.Info

	for _, p := range paths {
		_, err := os.Stat(p)
		if os.IsNotExist(err) {
			s.logger.Debugf("File (%v) does not exist.", p)
			continue
		} else if err != nil {
			return nil, fmt.Errorf("failed to check file: %v", err)
		}

		var output []byte
		if output, err = s.executeSSHKeyGenCommand("sha256", p); err != nil {
			return nil, fmt.Errorf("failed to execute ssh-keygen command: %v", err)
		}

		infos = append(infos, parseSSHKeyGenCommandOutput(string(output), infoType, p)...)
	}

	return infos, nil
}

func (s *Scanner) getPrivateKeysPaths(rootPath string, recursive bool) ([]string, error) {
	var paths []string
	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if path != rootPath && !recursive {
				return filepath.SkipDir
			}
			return nil
		}

		isPrivateKeyFile, err := isPrivateKey(path)
		if err != nil {
			s.logger.Errorf("failed to verify if file (%v) is private key file - skipping: %v", path, err)
			return nil
		}

		if isPrivateKeyFile {
			paths = append(paths, path)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return paths, nil
}

func isPrivateKey(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, fmt.Errorf("failed to open file: %v", err)
	}
	defer f.Close()

	// Splits on newlines by default.
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		// We only need to look at the first line
		return strings.Contains(scanner.Text(), "PRIVATE KEY"), nil
	}

	if err = scanner.Err(); err != nil {
		return false, fmt.Errorf("failed to scan file: %v", err)
	}

	return false, nil
}

func parseSSHKeyGenCommandOutput(output string, infoType types.InfoType, path string) []types.Info {
	var infos []types.Info
	lines := strings.Split(output, "\n")
	for i := range lines {
		if lines[i] == "" {
			continue
		}
		infos = append(infos, types.Info{
			Type: infoType,
			Path: path,
			Data: lines[i],
		})
	}
	return infos
}

func (s *Scanner) executeSSHKeyGenCommand(hashAlgo string, filePath string) ([]byte, error) {
	args := []string{
		"-E",
		hashAlgo,
		"-l",
		"-f",
		filePath,
	}
	cmd := exec.Command("ssh-keygen", args...)
	s.logger.Infof("Running command: %v", cmd.String())
	return sharedUtils.RunCommand(cmd)
}

func (s *Scanner) isValidInputType(sourceType utils.SourceType) bool {
	switch sourceType {
	case utils.ROOTFS:
		return true
	case utils.DIR, utils.FILE, utils.IMAGE, utils.SBOM:
		s.logger.Infof("source type %v is not supported for sshTopology, skipping.", sourceType)
	default:
		s.logger.Infof("unknown source type %v, skipping.", sourceType)
	}
	return false
}

func (s *Scanner) sendResults(results types.ScannerResult, err error) {
	if err != nil {
		s.logger.Error(err)
		results.Error = err
	}
	select {
	case s.resultChan <- &results:
	default:
		s.logger.Error("Failed to send results on channel")
	}
}
