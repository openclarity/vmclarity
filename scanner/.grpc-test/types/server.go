package types

import (
	"context"
	"host/proto"
)

type Scanner interface {
	GetInfo(ctx context.Context, empty *proto.Empty) (*proto.ScannerInfo, error)
	Scan(ctx context.Context, request *proto.ScanRequest) (*proto.ScanResult, error)
}
