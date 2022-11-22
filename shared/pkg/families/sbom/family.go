package sbom

import (
	log "github.com/sirupsen/logrus"

	_interface "github.com/openclarity/vmclarity/shared/pkg/families/interface"
)

type SBOM struct {
	conf Config
}

func (S SBOM) Run(getter _interface.ResultsGetter) (_interface.IsResults, error) {
	//TODO implement me
	log.Info("SBOM Run...")
	return &Results{}, nil
}

// ensure types implement the requisite interfaces
var _ _interface.Family = &SBOM{}

func New(conf Config) *SBOM {
	return &SBOM{
		conf: conf,
	}
}
