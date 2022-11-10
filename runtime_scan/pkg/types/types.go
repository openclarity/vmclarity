package types

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
	InstancesToScan          uint32
	InstancesStartedToScan   uint32
	InstancesCompletedToScan uint32
	Status                   Status
}

func (s *ScanProgress) SetStatus(status Status) {
	s.Status = status
}

type InstanceScanResult struct {
	// Instance data
	Instances Instance
	// Scan results
	Vulnerabilities []string // TODO define vulnerabilities struct
	Success         bool
	ScanErrors      []*ScanError
}

type ScanResults struct {
	InstanceScanResults []*InstanceScanResult
	Progress            ScanProgress
}

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
