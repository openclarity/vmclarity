package rootkits

import (
	log "github.com/sirupsen/logrus"

	_interface "github.com/openclarity/vmclarity/shared/pkg/families/interface"
)

type Rootkits struct {
	conf Config
}

func (r Rootkits) Run(getter _interface.ResultsGetter) (_interface.IsResults, error) {
	//TODO implement me
	log.Info("Rootkits Run...")
	return &Results{}, nil
}

// ensure types implement the requisite interfaces
var _ _interface.Family = &Rootkits{}

func New(conf Config) *Rootkits {
	return &Rootkits{
		conf: conf,
	}
}
