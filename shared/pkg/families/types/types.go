package types

type FamilyType string

const (
	SBOM FamilyType = "sbom"

	Vulnerabilities  FamilyType = "vulnerabilities"
	Secrets          FamilyType = "secrets"
	Rootkits         FamilyType = "rootkits"
	Malware          FamilyType = "malware"
	Misconfiguration FamilyType = "misconfiguration"

	Exploits FamilyType = "exploits"
)
