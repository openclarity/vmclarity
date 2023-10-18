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

package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/openclarity/vmclarity/pkg/shared/log"
)

func ReadFile(ctx context.Context, path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.GetLoggerFromContextOrDiscard(ctx).Errorf("failed to close file: %s", err.Error())
		}
	}()

	fileInfo, err := file.Stat()
	if err != nil {
		return []byte{}, fmt.Errorf("failed to get file info: %w", err)
	}

	fileSize := fileInfo.Size()
	buffer := make([]byte, fileSize)

	_, err = file.Read(buffer)
	if err != nil {
		return []byte{}, fmt.Errorf("failed to read file: %w", err)
	}

	return buffer, nil
}

func WriteFile(ctx context.Context, path string, data []byte) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer func() {
		err := file.Close()
		if err != nil {
			log.GetLoggerFromContextOrDiscard(ctx).Errorf("failed to close file: %s", err.Error())
		}
	}()

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write to file: %w", err)
	}

	return nil
}
