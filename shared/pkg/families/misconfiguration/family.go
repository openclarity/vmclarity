package misconfiguration

import (
	log "github.com/sirupsen/logrus"

	_interface "github.com/openclarity/vmclarity/shared/pkg/families/interface"
)

type Misconfiguration struct {
	conf Config
}

func (m Misconfiguration) Run(getter _interface.ResultsGetter) (_interface.IsResults, error) {
	//TODO implement me
	log.Info("Misconfiguration Run...")
	return &Results{}, nil
}

// ensure types implement the requisite interfaces
var _ _interface.Family = &Misconfiguration{}

func New(conf Config) *Misconfiguration {
	return &Misconfiguration{
		conf: conf,
	}
}
