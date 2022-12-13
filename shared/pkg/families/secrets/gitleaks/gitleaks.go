package gitleaks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/openclarity/kubeclarity/shared/pkg/job_manager"
	"github.com/openclarity/kubeclarity/shared/pkg/utils"
	"github.com/openclarity/vmclarity/shared/pkg/families/secrets"
	"github.com/openclarity/vmclarity/shared/pkg/families/secrets/common"
	log "github.com/sirupsen/logrus"
)

const ScannerName = "gitleaks"
const reportPath = "/tmp/gitleaks.json"

type Scanner struct {
	name       string
	logger     *log.Entry
	config     Config
	resultChan chan job_manager.Result
}

func New(c job_manager.IsConfig, logger *log.Entry, resultChan chan job_manager.Result) job_manager.Job {
	conf := c.(*secrets.ScannersConfig) // nolint:forcetypeassert
	return &Scanner{
		name:       ScannerName,
		logger:     logger.Dup().WithField("secret-scanner", ScannerName),
		config:     Config{BinaryPath: conf.Gitleaks.BinaryPath},
		resultChan: resultChan,
	}
}

func (a *Scanner) Run(sourceType utils.SourceType, userInput string) error {
	go func() {
		retResults := common.Results{
			Source:      userInput,
			ScannerName: ScannerName,
		}
		// validate that gitleaks binary exists
		if _, err := os.Stat(a.config.BinaryPath); err != nil {
			retResults.Error = fmt.Errorf("failed to find binary in %v: %v", a.config.BinaryPath, err)
			a.resultChan <- &retResults
			return
		}

		// ./gitleaks detect --source=<source> --no-git -r <report-path> -f json --exit-code 0
		cmd := exec.Command(a.config.BinaryPath, "detect", fmt.Sprintf("--source=%v", userInput), "--no-git", "-r", reportPath, "-f", "json", "--exit-code", "0")
		_, err := runCommand(cmd)
		if err != nil {
			retResults.Error = fmt.Errorf("failed to run gitleaks command: %v", err)
			a.resultChan <- &retResults
			return
		}

		out, err := os.ReadFile(reportPath)
		if err != nil {
			retResults.Error = fmt.Errorf("failed to read report file from path: %v. %v", reportPath, err)
			a.resultChan <- &retResults
			return
		}

		if err := json.Unmarshal(out, &retResults.Findings); err != nil {
			retResults.Error = fmt.Errorf("failed to unmarshal results. out: %s. err: %v", out, err)
			a.resultChan <- &retResults
			return
		}
		a.resultChan <- &retResults
	}()

	return nil
}

func runCommand(cmd *exec.Cmd) ([]byte, error) {
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to run command %v, error: %w, stdout: %v, stderr: %v", cmd.String(), err, outb.String(), errb.String())
	}
	return outb.Bytes(), nil
}
