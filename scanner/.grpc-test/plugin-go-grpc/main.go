// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"google.golang.org/grpc"
	"host/proto" // versioned
	"host/types"
	"log"
	"net"
	"time"
)

type KV struct{}

func (k KV) GetInfo(ctx context.Context, empty *proto.Empty) (*proto.ScannerInfo, error) {
	return &proto.ScannerInfo{
		Metadata: &proto.Metadata{
			Annotations: map[string]string{
				"runtime": "unix",
			},
		},
		Name:    "example",
		Version: "1.2.3",
	}, nil
}

func (k KV) Scan(ctx context.Context, request *proto.ScanRequest) (*proto.ScanResult, error) {
	time.Sleep(10 * time.Second)
	return &proto.ScanResult{
		Findings: []*proto.Finding{
			{
				FindingInfo: &proto.Finding_Package{
					Package: &proto.Package{
						Name:     "package",
						Version:  "version",
						Type:     "type",
						Language: "",
						Licenses: nil,
						Cpes:     nil,
						Purl:     "something",
					},
				},
			},
		},
		Error:   "",
		Summary: nil,
	}, nil
}

func main() {
	grpcServer := grpc.NewServer()
	kvServer := &KV{}
	types.RegisterScanner(grpcServer, kvServer)

	listen, err := net.Listen("tcp", "0.0.0.0:3000")
	if err != nil {
		log.Fatalf("could not listen to 0.0.0.0:3000 %v", err)
	}

	log.Println("Server starting...")
	log.Fatal(grpcServer.Serve(listen))
}
