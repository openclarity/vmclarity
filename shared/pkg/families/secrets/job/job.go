package job

import (
	"github.com/openclarity/kubeclarity/shared/pkg/job_manager"
	"github.com/openclarity/vmclarity/shared/pkg/families/secrets/gitleaks"
)

var Factory = job_manager.NewJobFactory()

func init() {
	Factory.Register(gitleaks.ScannerName, gitleaks.New)
}
