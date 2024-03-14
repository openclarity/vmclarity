// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"host/types"
	"os"
)

func run() error {
	// We're a host. Start by launching the plugin process.
	conn, err := grpc.Dial("0.0.0.0:3000", grpc.WithCredentialsBundle(insecure.NewBundle()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{}))
	if err != nil {
		return err
	}
	defer conn.Close()

	// Request the plugin
	raw, err := types.NewScannerClient(conn)
	if err != nil {
		return err
	}

	// We should have a KV store now! This feels like a normal interface
	// implementation but is in fact over an RPC connection.
	os.Args = os.Args[1:]
	switch os.Args[0] {
	case "get":
		result, err := raw.GetInfo(context.Background(), nil)
		if err != nil {
			return err
		}

		fmt.Println("scanner-info", result)

	case "scan":
		result, err := raw.Scan(context.Background(), nil)
		if err != nil {
			return err
		}

		fmt.Println("scan-result", result)

	default:
		return fmt.Errorf("Please only use 'get' or 'put', given: %q", os.Args[0])
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Printf("error: %+v\n", err)
		os.Exit(1)
	}

	os.Exit(0)
}
