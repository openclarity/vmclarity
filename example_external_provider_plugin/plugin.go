package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	provider_service "github.com/openclarity/vmclarity/runtime_scan/pkg/provider/external/proto"
)

type Provider struct {
	provider_service.UnimplementedProviderServer
	ReadyToScan bool
}

func (p *Provider) DiscoverAssets(_ context.Context, _ *provider_service.DiscoverAssetsParams) (*provider_service.DiscoverAssetsResult, error) {
	return &provider_service.DiscoverAssetsResult{
		Assets: []*provider_service.Asset{
			{
				AssetType: &provider_service.Asset_Vminfo{
					Vminfo: &provider_service.VMInfo{
						Id:           "Id",
						Location:     "Location",
						Image:        "Image",
						InstanceType: "InstanceType",
						Platform:     "Platform",
						Tags: []*provider_service.Tag{
							{
								Key: "key",
								Val: "val",
							},
						},
						LaunchTime: timestamppb.New(time.Now()),
					},
				},
			},
		},
	}, nil
}

func (p *Provider) RemoveAssetScan(context.Context, *provider_service.RemoveAssetScanParams) (*provider_service.RemoveAssetScanResult, error) {
	// clean all resources that were created during RunAssetScan
	return &provider_service.RemoveAssetScanResult{}, nil
}

func (p *Provider) RunAssetScan(context.Context, *provider_service.RunAssetScanParams) (*provider_service.RunAssetScanResult, error) {
	// Create all resources needed in order to start the scan.
	// It can be spinning up a VM or snapshoting a volume.
	// It should be non blocking and idompetent...

	if !p.ReadyToScan {
		// flip the ready to scan bit, so next time we're getting called, ERR_NONE will be returned.
		p.ReadyToScan = true

		// Tell VMClarity that resources are not ready and RunAssetScan should be called again in the next iteration.
		return &provider_service.RunAssetScanResult{
			ErrType: provider_service.ErrorType_ERR_RETRYABLE,
		}, fmt.Errorf("not all resources are ready")
	}

	// when all the resource creation is done and ready, an error type of ErrorType_ERR_NONE should be return.
	return &provider_service.RunAssetScanResult{
		ErrType: provider_service.ErrorType_ERR_NONE,
	}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:24230"))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	provider_service.RegisterProviderServer(grpcServer, &Provider{})
	log.Infof("listening.....")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
