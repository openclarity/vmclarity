package types

type Tag struct {
	Key string
	Val string
}

type SecurityGroup struct {
	ID string
}

type VPC struct {
	ID             string
	SecurityGroups []SecurityGroup
}

type Region struct {
	ID   string
	VPCs []VPC
}

type ScanScope struct {
	All         bool
	Regions     []Region
	ScanStopped bool
	IncludeTags []*Tag
	ExcludeTags []*Tag
}

type Status string

const (
	Idle            Status = "Idle"
	ScanInit        Status = "ScanInit"
	ScanInitFailure Status = "ScanInitFailure"
	NothingToScan   Status = "NothingToScan"
	Scanning        Status = "Scanning"
	DoneScanning    Status = "DoneScanning"
)

type ScanProgress struct {
	ImagesToScan          uint32
	ImagesStartedToScan   uint32
	ImagesCompletedToScan uint32
	Status                Status
}

func (s *ScanProgress) SetStatus(status Status) {
	s.Status = status
}

type VMScanResult struct {
	// VM data
	VMID string
	// Scan results
	Vulnerabilities []string
	Success         bool
}

type ScanResults struct {
	ImageScanResults []*VMScanResult
	Progress         ScanProgress
}

type Job struct {
	Instance    Instance
	SrcSnapshot Snapshot
	DstSnapshot Snapshot
}

type Instance struct {
	ID     string
	Region string
}

type Snapshot struct {
	ID     string
	Region string
}

type Volume struct {
	ID     string
	Name   string
	Region string
}
