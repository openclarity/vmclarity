package secrets

import (
	"github.com/openclarity/vmclarity/shared/pkg/families/secrets/common"
)

type MergedResults struct {
	Results []*common.Results
}

func NewMergedResults() *MergedResults {
	return &MergedResults{}
}

func (m *MergedResults) Merge(other *common.Results) *MergedResults {
	m.Results = append(m.Results, other)
	return m
}
