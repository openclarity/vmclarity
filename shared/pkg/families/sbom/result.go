package sbom

import cdx "github.com/CycloneDX/cyclonedx-go"

type Results struct {
	BOM cdx.BOM
}

func (*Results) IsResults() {}
