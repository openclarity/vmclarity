package aws

type ScanScope struct {
	All         bool
	Regions     []Region
	ScanStopped bool
	IncludeTags []Tag
	ExcludeTags []Tag
}

type Tag struct {
	key string
	val string
}

type SecurityGroup struct {
	id string
}

type VPC struct {
	Id             string
	securityGroups []SecurityGroup
}

type Region struct {
	Id   string
	vpcs []VPC
}
