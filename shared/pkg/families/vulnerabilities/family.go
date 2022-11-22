package vulnerabilities

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	_interface "github.com/openclarity/vmclarity/shared/pkg/families/interface"
	"github.com/openclarity/vmclarity/shared/pkg/families/sbom"
	"github.com/openclarity/vmclarity/shared/pkg/families/types"
)

type Vulnerabilities struct {
	conf Config
}

func (v Vulnerabilities) Run(getter _interface.ResultsGetter) (_interface.IsResults, error) {
	//TODO implement me
	log.Info("Vulnerabilities Run...")

	// Example for InputFromFamily
	for _, family := range v.conf.InputFromFamily {
		results := getter.GetResults(family.FamilyType)
		switch family.FamilyType {
		case types.SBOM:
			sbomResults, ok := results.(*sbom.Results)
			if !ok {
				return nil, fmt.Errorf("failed to cast sbom results")
			}
			log.Info("SBOM results: %+v", sbomResults)
		default:
			return nil, fmt.Errorf("not supported family as input: %v", family.FamilyType)
		}
	}

	return &Results{}, nil
}

// ensure types implement the requisite interfaces
var _ _interface.Family = &Vulnerabilities{}

func New(conf Config) *Vulnerabilities {
	return &Vulnerabilities{
		conf: conf,
	}
}
