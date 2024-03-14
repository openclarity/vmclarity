// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

// Package shared contains shared data between the host and plugins.
package types

import (
	"google.golang.org/grpc"

	"host/proto"
)

func RegisterScanner(s *grpc.Server, scanner Scanner) {
	proto.RegisterScannerServer(s, scanner)
}

func NewScannerClient(c *grpc.ClientConn) (proto.ScannerClient, error) {
	return proto.NewScannerClient(c), nil
}
