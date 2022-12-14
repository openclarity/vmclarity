package gitleaks

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	sharedutils "github.com/openclarity/vmclarity/shared/pkg/utils"
	log "github.com/sirupsen/logrus"

	"github.com/openclarity/kubeclarity/shared/pkg/job_manager"
	"github.com/openclarity/kubeclarity/shared/pkg/utils"
	"github.com/openclarity/vmclarity/shared/pkg/families/secrets/common"
	gitleaksconfig "github.com/openclarity/vmclarity/shared/pkg/families/secrets/gitleaks/config"
)

const ScannerName = "gitleaks"
const reportPath = "/tmp/gitleaks.json"

type Scanner struct {
	name       string
	logger     *log.Entry
	config     gitleaksconfig.Config
	resultChan chan job_manager.Result
}

func New(c job_manager.IsConfig, logger *log.Entry, resultChan chan job_manager.Result) job_manager.Job {
	conf := c.(*common.ScannersConfig) // nolint:forcetypeassert
	return &Scanner{
		name:       ScannerName,
		logger:     logger.Dup().WithField("scanner", ScannerName),
		config:     gitleaksconfig.Config{BinaryPath: conf.Gitleaks.BinaryPath},
		resultChan: resultChan,
	}
}

func (a *Scanner) Run(sourceType utils.SourceType, userInput string) error {
	if sourceType != utils.DIR {
		return fmt.Errorf("invalid source type for gitleaks: %v", sourceType)
	}
	go func() {
		retResults := common.Results{
			Source:      userInput,
			ScannerName: ScannerName,
		}
		// validate that gitleaks binary exists
		if _, err := os.Stat(a.config.BinaryPath); err != nil {
			a.sendResults(retResults, fmt.Errorf("failed to find binary in %v: %v", a.config.BinaryPath, err))
			return
		}

		// ./gitleaks detect --source=<source> --no-git -r <report-path> -f json --exit-code 0
		cmd := exec.Command(a.config.BinaryPath, "detect", fmt.Sprintf("--source=%v", userInput), "--no-git", "-r", reportPath, "-f", "json", "--exit-code", "0")
		a.logger.Infof("running gitleaks command: %v", cmd.String())
		_, err := sharedutils.RunCommand(cmd)
		if err != nil {
			a.sendResults(retResults, fmt.Errorf("failed to run gitleaks command: %v", err))
			return
		}

		out, err := os.ReadFile(reportPath)
		if err != nil {
			a.sendResults(retResults, fmt.Errorf("failed to read report file from path %v: %v", reportPath, err))
			return
		}
		defer func() {
			_ = os.Remove(reportPath)
		}()

		if err := json.Unmarshal(out, &retResults.Findings); err != nil {
			a.sendResults(retResults, fmt.Errorf("failed to unmarshal results. out: %s. err: %v", out, err))
			return
		}
		a.sendResults(retResults, nil)
	}()

	return nil
}

func (a *Scanner) sendResults(results common.Results, err error) {
	if err != nil {
		a.logger.Error(err)
		results.Error = err
	}
	select {
	case a.resultChan <- &results:
	default:
		a.logger.Error("Failed to send results on channel")
	}
}
