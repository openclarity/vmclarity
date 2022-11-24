package vulnerabilities

import (
	"fmt"
	"os"

	"github.com/openclarity/kubeclarity/shared/pkg/config"
	"github.com/openclarity/kubeclarity/shared/pkg/scanner/job"
	"github.com/openclarity/kubeclarity/shared/pkg/utils"
	log "github.com/sirupsen/logrus"

	"github.com/openclarity/kubeclarity/shared/pkg/job_manager"
	sharedscanner "github.com/openclarity/kubeclarity/shared/pkg/scanner"

	_interface "github.com/openclarity/vmclarity/shared/pkg/families/interface"
	"github.com/openclarity/vmclarity/shared/pkg/families/sbom"
	"github.com/openclarity/vmclarity/shared/pkg/families/types"
)

const (
	sbomTempFilePath = "/tmp/sbom"
)

type Vulnerabilities struct {
	logger         *log.Entry
	conf           Config
	ScannersConfig config.Config
}

func (v Vulnerabilities) Run(getter _interface.ResultsGetter) (_interface.IsResults, error) {
	v.logger.Info("Vulnerabilities Run...")

	manager := job_manager.New(v.conf.ScannersList, v.conf.ScannersConfig, v.logger, job.CreateJob)
	mergedResults := sharedscanner.NewMergedResults()

	if v.conf.InputFromSbom {
		results := getter.GetResults(types.SBOM)
		sbomResults, ok := results.(*sbom.Results)
		if !ok {
			return nil, fmt.Errorf("failed to cast sbom results")
		}
		v.logger.Infof("Using input from SBOM results")

		// TODO: need to avoid writing sbom to file
		//
		//dx := formatter.New(sbomResults.Format, sbomResults.SBOM)
		//err := dx.SetSBOM(sbomResults.BOM)
		//if err != nil {
		//	return nil, err
		//}
		//err = dx.Encode(formatter.CycloneDXFormat)
		//if err != nil {
		//	return nil, err
		//}

		if err := os.WriteFile(sbomTempFilePath, sbomResults.SBOM, 0600 /* read & write */); err != nil { // nolint:gomnd,gofumpt
			return nil, fmt.Errorf("failed to write sbom to file: %v", err)
		}

		v.conf.Inputs = append(v.conf.Inputs, Inputs{
			Input:     sbomTempFilePath,
			InputType: "sbom",
		})
	}

	for _, input := range v.conf.Inputs {
		results, err := manager.Run(utils.SourceType(input.InputType), input.Input)
		if err != nil {
			return nil, err
		}

		// Merge results.
		for name, result := range results {
			v.logger.Infof("Merging result from %q", name)
			mergedResults = mergedResults.Merge(result.(*sharedscanner.Results)) // nolint:forcetypeassert
		}

		// TODO:
		//// Set source values.
		//mergedResults.SetSource(sharedscanner.Source{
		//	Type: "image",
		//	Name: config.ImageIDToScan,
		//	Hash: config.ImageHashToScan,
		//})
	}

	v.logger.Info("Vulnerabilities Done...")

	return &Results{
		MergedResults: mergedResults,
	}, nil
}

// ensure types implement the requisite interfaces
var _ _interface.Family = &Vulnerabilities{}

func New(logger *log.Entry, conf Config) *Vulnerabilities {
	return &Vulnerabilities{
		logger: logger.Dup().WithField("family", "vulnerabilities"),
		conf:   conf,
	}
}
