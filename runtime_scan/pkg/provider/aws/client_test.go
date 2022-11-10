package aws

import (
	"context"
	"testing"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"

	"github.com/openclarity/vmclarity/runtime_scan/pkg/types"
)

func TestClient_ListAllRegions(t *testing.T) {
	cfg, err := awsconfig.LoadDefaultConfig(context.TODO())
	if err != nil {
		t.Fatalf("%v", err)
	}
	ec2Client := ec2.NewFromConfig(cfg)

	c := Client{
		ec2Client: ec2Client,
	}
	instance := types.Instance{
		ID:     "i-0f70b335ea12b2853",
		Region: "us-east-1",
	}
	rootVolume, err := c.GetInstanceRootVolume(instance)
	if err != nil {
		t.Fatalf("%v", err)
	}
	// create a snapshot of that vm
	srcSnapshot, err := c.CreateSnapshot(rootVolume)
	if err != nil {
		t.Fatalf("%v", err)
	}
	if err := c.WaitForSnapshotReady(srcSnapshot); err != nil {
		t.Fatalf("%v", err)
	}
	//copy the snapshot to the scanner region
	// TODO make sure we need this.
	cpySnapshot, err := c.CopySnapshot(srcSnapshot, "us-east-2")
	if err != nil {
		t.Fatalf("%v", err)
	}
	if err := c.WaitForSnapshotReady(cpySnapshot); err != nil {
		t.Fatalf("%v", err)
	}
	// create the scanner job (vm) with a boot script
	launchedInstance, err := c.LaunchInstance("ami-0568773882d492fc8", "xvdh", cpySnapshot)
	if err != nil {
		t.Fatalf("%v", err)
	}

	t.Logf("res: %v", launchedInstance.ID)
}
