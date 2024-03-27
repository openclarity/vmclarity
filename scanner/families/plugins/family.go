package plugins

import (
	"context"

	"github.com/openclarity/vmclarity/scanner/families/interfaces"
	"github.com/openclarity/vmclarity/scanner/families/results"
	"github.com/openclarity/vmclarity/scanner/families/types"
)

type Plugins struct {
	conf Config
}

var _ interfaces.Family = &Plugins{}

func (p *Plugins) Run(ctx context.Context, res *results.Results) (interfaces.IsResults, error) {
	return nil, nil
}

func (p *Plugins) GetType() types.FamilyType {
	return types.Plugins
}

func New(conf Config) *Plugins {
	return &Plugins{
		conf: conf,
	}
}
