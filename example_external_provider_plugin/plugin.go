package main

import (
	"context"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	provider_service "github.com/openclarity/vmclarity/runtime_scan/pkg/provider/grpc/proto"

	"github.com/openclarity/vmclarity/api/models"
)

type Provider struct {
	provider_service.UnimplementedProviderServer
}

func (p *Provider) RunScanningJob(ctx context.Context, params *provider_service.RunScanningJobParams) (*provider_service.RunScanningJobResult, error) {
	// Here need to run the VMClarity CLI and return the instance that is running the CLI.
	return &provider_service.RunScanningJobResult{
		Instance: &provider_service.Instance{
			Id:           params.Id,
			Location:     params.Location,
			Image:        "ami-01d08089481510ba2",
			InstanceType: "t2.large",
			Platform:     "Linux/UNIX",
			Tags: []*provider_service.Tag{
				{
					Key: "Name",
					Val: "cloud provider plugin instance",
				},
			},
			LaunchTime: timestamppb.New(time.Now()),
		},
	}, nil
}

func (p *Provider) DiscoverScopes(ctx context.Context, params *provider_service.DiscoverScopesParams) (*provider_service.DiscoverScopesResult, error) {
	// Scopes can be very cloud provider specific. Here you see an example of an AWS scopes which
	// consist of regions->vpcs->securityGroups.
	// in order to support a new scope, needs to add a new scopeType object in VMClarity backend api.
	scopes := models.ScopeType{}

	err := scopes.FromAwsAccountScope(models.AwsAccountScope{
		Regions: &[]models.AwsRegion{
			{
				Name: "region1",
				Vpcs: &[]models.AwsVPC{
					{
						Id: "vpc-1",
						SecurityGroups: &[]models.AwsSecurityGroup{
							{
								Id: "sg-1",
							},
						},
					},
					{
						Id: "vpc-2",
						SecurityGroups: &[]models.AwsSecurityGroup{
							{
								Id: "sg-2",
							},
						},
					},
				},
			},
			{
				Name: "region2",
				Vpcs: &[]models.AwsVPC{
					{
						Id: "vpc3",
						SecurityGroups: &[]models.AwsSecurityGroup{
							{
								Id: "sg3",
							},
						},
					},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	scopesB, err := scopes.MarshalJSON()
	if err != nil {
		return nil, err
	}

	return &provider_service.DiscoverScopesResult{
		Scopes: string(scopesB),
	}, nil
}

func (p *Provider) DiscoverInstances(context.Context, *provider_service.DiscoverInstancesParams) (*provider_service.DiscoverInstancesResult, error) {
	// Here needs to discover all instances according to the ScanScope parameter.
	// ScanScope can be added to the backend api, and can contain fields like instance tags, location etc.
	return &provider_service.DiscoverInstancesResult{
		Instances: []*provider_service.Instance{
			{
				Id:           "instanceID1",
				Location:     "region1",
				Image:        "ami-01d08089481510ba2",
				InstanceType: "t2.large",
				Platform:     "Linux/UNIX",
				Tags: []*provider_service.Tag{
					{
						Key: "Name",
						Val: "instance1",
					},
				},
				LaunchTime: timestamppb.New(time.Now()),
			},
			{
				Id:           "instanceID2",
				Location:     "region2",
				Image:        "ami-01d08089481510ba2",
				InstanceType: "t2.large",
				Platform:     "Linux/UNIX",
				Tags: []*provider_service.Tag{
					{
						Key: "Name",
						Val: "instance2",
					},
				},
				LaunchTime: timestamppb.New(time.Now()),
			},
		},
	}, nil
}

func (p *Provider) GetInstanceRootVolume(context.Context, *provider_service.GetInstanceRootVolumeParams) (*provider_service.GetInstanceRootVolumeResult, error) {
	return &provider_service.GetInstanceRootVolumeResult{
		Volume: &provider_service.Volume{
			Id:       "volume1",
			Location: "region1",
		},
	}, nil
}

func (p *Provider) WaitForInstanceReady(context.Context, *provider_service.WaitForInstanceReadyParams) (*provider_service.WaitForInstanceReadyResult, error) {
	return &provider_service.WaitForInstanceReadyResult{}, nil
}

func (p *Provider) DeleteInstance(context.Context, *provider_service.DeleteInstanceParams) (*provider_service.DeleteInstanceResult, error) {
	return &provider_service.DeleteInstanceResult{}, nil
}

func (p *Provider) AttachVolumeToInstance(context.Context, *provider_service.AttachVolumeToInstanceParams) (*provider_service.AttachVolumeToInstanceResult, error) {
	return &provider_service.AttachVolumeToInstanceResult{}, nil
}

func (p *Provider) CopySnapshot(context.Context, *provider_service.CopySnapshotParams) (*provider_service.CopySnapshotResult, error) {
	return &provider_service.CopySnapshotResult{
		Snapshot: &provider_service.Snapshot{
			Id:       "snapshot1",
			Location: "region2",
		},
	}, nil
}

func (p *Provider) DeleteSnapshot(context.Context, *provider_service.DeleteSnapshotParams) (*provider_service.DeleteSnapshotResult, error) {
	return &provider_service.DeleteSnapshotResult{}, nil
}

func (p *Provider) WaitForSnapshotReady(context.Context, *provider_service.WaitForSnapshotReadyParams) (*provider_service.WaitForSnapshotReadyResult, error) {
	return &provider_service.WaitForSnapshotReadyResult{}, nil
}

func (p *Provider) CreateVolumeFromSnapshot(context.Context, *provider_service.CreateVolumeFromSnapshotParams) (*provider_service.CreateVolumeFromSnapshotResult, error) {
	return &provider_service.CreateVolumeFromSnapshotResult{
		Volume: &provider_service.Volume{
			Id:       "volume2",
			Location: "region2",
		},
	}, nil
}

func (p *Provider) TakeVolumeSnapshot(context.Context, *provider_service.TakeVolumeSnapshotParams) (*provider_service.TakeVolumeSnapshotResult, error) {
	return &provider_service.TakeVolumeSnapshotResult{
		Snapshot: &provider_service.Snapshot{
			Id:       "snapshot11",
			Location: "region1",
		},
	}, nil
}

func (p *Provider) WaitForVolumeReady(context.Context, *provider_service.WaitForVolumeReadyParams) (*provider_service.WaitForVolumeReadyResult, error) {
	return &provider_service.WaitForVolumeReadyResult{}, nil
}

func (p *Provider) WaitForVolumeAttached(context.Context, *provider_service.WaitForVolumeAttachedParams) (*provider_service.WaitForVolumeAttachedResult, error) {
	return &provider_service.WaitForVolumeAttachedResult{}, nil
}

func (p *Provider) DeleteVolume(context.Context, *provider_service.DeleteVolumeParams) (*provider_service.DeleteVolumeResult, error) {
	return &provider_service.DeleteVolumeResult{}, nil
}

//
//func main() {
//	flag.Parse()
//	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:24230"))
//	if err != nil {
//		log.Fatalf("failed to listen: %v", err)
//	}
//	var opts []grpc.ServerOption
//	grpcServer := grpc.NewServer(opts...)
//	provider_service.RegisterProviderServer(grpcServer, &Provider{})
//	log.Infof("listening.....")
//	if err := grpcServer.Serve(lis); err != nil {
//		log.Fatalf("failed to serve: %v", err)
//	}
//}
