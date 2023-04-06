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

package presenter

import (
	"context"
	"fmt"

	"github.com/openclarity/vmclarity/shared/pkg/families"
	"github.com/openclarity/vmclarity/shared/pkg/families/results"
)

type MultiPresenter struct {
	Presenters []Presenter
}

func (m *MultiPresenter) ExportSbomResult(ctx context.Context, res *results.Results, famErr families.RunErrors) error {
	for _, p := range m.Presenters {
		if err := p.ExportSbomResult(ctx, res, famErr); err != nil {
			return fmt.Errorf("failed to export result: %w", err)
		}
	}

	return nil
}

func (m *MultiPresenter) ExportVulResult(ctx context.Context, res *results.Results, famErr families.RunErrors) error {
	for _, p := range m.Presenters {
		if err := p.ExportVulResult(ctx, res, famErr); err != nil {
			return fmt.Errorf("failed to export result: %w", err)
		}
	}

	return nil
}

func (m *MultiPresenter) ExportSecretsResult(ctx context.Context, res *results.Results, famErr families.RunErrors) error {
	for _, p := range m.Presenters {
		if err := p.ExportSecretsResult(ctx, res, famErr); err != nil {
			return fmt.Errorf("failed to export result: %w", err)
		}
	}

	return nil
}

func (m *MultiPresenter) ExportMalwareResult(ctx context.Context, res *results.Results, famErr families.RunErrors) error {
	for _, p := range m.Presenters {
		if err := p.ExportMalwareResult(ctx, res, famErr); err != nil {
			return fmt.Errorf("failed to export result: %w", err)
		}
	}

	return nil
}

func (m *MultiPresenter) ExportExploitsResult(ctx context.Context, res *results.Results, famErr families.RunErrors) error {
	for _, p := range m.Presenters {
		if err := p.ExportExploitsResult(ctx, res, famErr); err != nil {
			return fmt.Errorf("failed to export result: %w", err)
		}
	}

	return nil
}

func (m *MultiPresenter) ExportMisconfigurationResult(ctx context.Context, res *results.Results, famErr families.RunErrors) error {
	for _, p := range m.Presenters {
		if err := p.ExportMisconfigurationResult(ctx, res, famErr); err != nil {
			return fmt.Errorf("failed to export result: %w", err)
		}
	}

	return nil
}
