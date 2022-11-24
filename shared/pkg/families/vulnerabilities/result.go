package vulnerabilities

import "github.com/openclarity/kubeclarity/shared/pkg/scanner"

type Results struct {
	MergedResults *scanner.MergedResults
}

func (*Results) IsResults() {}
