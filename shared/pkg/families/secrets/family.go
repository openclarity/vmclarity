package secrets

import (
	log "github.com/sirupsen/logrus"

	_interface "github.com/openclarity/vmclarity/shared/pkg/families/interface"
)

type Secrets struct {
	conf Config
}

func (s Secrets) Run(getter _interface.ResultsGetter) (_interface.IsResults, error) {
	//TODO implement me
	log.Info("Secrets Run...")
	return &Results{}, nil
}

// ensure types implement the requisite interfaces
var _ _interface.Family = &Secrets{}

func New(conf Config) *Secrets {
	return &Secrets{
		conf: conf,
	}
}
