package sbom

import (
	"fmt"

	"github.com/openclarity/kubeclarity/cli/pkg"
	cliutils "github.com/openclarity/kubeclarity/cli/pkg/utils"
	sharedanalyzer "github.com/openclarity/kubeclarity/shared/pkg/analyzer"
	"github.com/openclarity/kubeclarity/shared/pkg/analyzer/job"
	"github.com/openclarity/kubeclarity/shared/pkg/job_manager"
	"github.com/openclarity/kubeclarity/shared/pkg/utils"
	log "github.com/sirupsen/logrus"

	_interface "github.com/openclarity/vmclarity/shared/pkg/families/interface"
)

type SBOM struct {
	logger *log.Entry
	conf   Config
}

func (s SBOM) Run(_ _interface.ResultsGetter) (_interface.IsResults, error) {
	s.logger.Info("SBOM Run...")

	if len(s.conf.Inputs) == 0 {
		return nil, fmt.Errorf("inputs list is empty")
	}

	outputFormat := s.conf.AnalyzersConfig.Analyzer.OutputFormat

	// TODO: move the logic from cli utils to shared utils
	// TODO: now that we support multiple inputs,
	//  we need to change the fact the the MergedResults assumes it is only for 1 input?
	hash, err := cliutils.GenerateHash(utils.SourceType(s.conf.Inputs[0].InputType), s.conf.Inputs[0].Input)
	if err != nil {
		return nil, fmt.Errorf("failed to generate hash for source %s: %v", s.conf.Inputs[0].Input, err)
	}

	manager := job_manager.New(s.conf.AnalyzersList, s.conf.AnalyzersConfig, s.logger, job.CreateAnalyzerJob)
	mergedResults := sharedanalyzer.NewMergedResults(utils.SourceType(s.conf.Inputs[0].InputType), hash)

	for _, input := range s.conf.Inputs {
		results, err := manager.Run(utils.SourceType(input.InputType), input.Input)
		if err != nil {
			return nil, err
		}

		// Merge results.
		for name, result := range results {
			s.logger.Infof("Merging result from %q", name)
			mergedResults = mergedResults.Merge(result.(*sharedanalyzer.Results), outputFormat) // nolint:forcetypeassert
		}
	}

	MergeWithResults := make(map[string]job_manager.Result)
	for i, with := range s.conf.MergeWith {
		name := fmt.Sprintf("merge_with_%d", i)
		cdxBOMBytes, err := cliutils.ConvertInputSBOMIfNeeded(with.SbomPath, outputFormat)
		if err != nil {
			return nil, fmt.Errorf("failed to convert merged with SBOM. path=%s: %v", with.SbomPath, err)
		}
		MergeWithResults[name] = sharedanalyzer.CreateResults(cdxBOMBytes, name, with.SbomPath, utils.SBOM)
	}

	mergedSBOMBytes, err := mergedResults.CreateMergedSBOMBytes(outputFormat, pkg.GitRevision)
	if err != nil {
		return nil, fmt.Errorf("failed to create merged output: %v", err)
	}

	s.logger.Info("SBOM Done...")

	return &Results{
		Format: outputFormat,
		SBOM:   mergedSBOMBytes,
	}, nil
}

// ensure types implement the requisite interfaces
var _ _interface.Family = &SBOM{}

func New(logger *log.Entry, conf Config) *SBOM {
	return &SBOM{
		conf:   conf,
		logger: logger.Dup().WithField("family", "sbom"),
	}
}
