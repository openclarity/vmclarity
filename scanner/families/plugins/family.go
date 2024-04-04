package plugins

import (
	"context"
	"fmt"
	"time"

	"github.com/openclarity/vmclarity/core/log"
	"github.com/openclarity/vmclarity/scanner/families/interfaces"
	"github.com/openclarity/vmclarity/scanner/families/plugins/job"
	"github.com/openclarity/vmclarity/scanner/families/results"
	"github.com/openclarity/vmclarity/scanner/families/types"
	familiesutils "github.com/openclarity/vmclarity/scanner/families/utils"
	"github.com/openclarity/vmclarity/scanner/job_manager"
	"github.com/openclarity/vmclarity/scanner/utils"
)

type Plugins struct {
	conf Config
}

var _ interfaces.Family = &Plugins{}

func (p *Plugins) Run(ctx context.Context, res *results.Results) (interfaces.IsResults, error) {
	logger := log.GetLoggerFromContextOrDiscard(ctx).WithField("family", "plugins")
	logger.Info("Plugins Run...")

	manager := job_manager.New(p.conf.ScannersList, p.conf.ScannersConfig, logger, job.Factory)

	var pluginsResults Results
	for _, input := range p.conf.Inputs {
		startTime := time.Now()
		_, err := manager.Run(utils.SourceType(input.InputType), input.Input)
		if err != nil {
			return nil, fmt.Errorf("failed to scan input %q for plugins: %w", input.Input, err)
		}
		endTime := time.Now()
		inputSize, err := familiesutils.GetInputSize(input)
		if err != nil {
			logger.Warnf("Failed to calculate input %v size: %v", input, err)
		}

		// TODO Add results to pluginsResults
		_ = types.CreateInputScanMetadata(startTime, endTime, inputSize, input)
	}

	logger.Info("Plugins Done...")
	return &pluginsResults, nil
}

func (p *Plugins) GetType() types.FamilyType {
	return types.Plugins
}

func New(conf Config) *Plugins {
	return &Plugins{
		conf: conf,
	}
}
