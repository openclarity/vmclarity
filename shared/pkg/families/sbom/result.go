package sbom

type Results struct {
	// TODO: Do we want to keep it cdx or []byte & Format?
	Format string
	SBOM   []byte
	//BOM cdx.BOM
}

func (*Results) IsResults() {}
