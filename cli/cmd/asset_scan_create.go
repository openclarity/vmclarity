/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/shared/pkg/backendclient"
	"github.com/openclarity/vmclarity/shared/pkg/utils"
)

// standaloneCmd represents the standalone command
var assetScanCreateCmd = &cobra.Command{
	Use:   "asset-scan-create",
	Short: "Create asset scan",
	Long:  `It creates asset scan. It's useful in the CI/CD mode without VMClarity orchestration`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Infof("asset-scan-create called")
		assetID, err := cmd.Flags().GetString("asset-id")
		if err != nil {
			logger.Fatalf("Unable to get asset id: %v", err)
		}
		server, err := cmd.Flags().GetString("server")
		if err != nil {
			logger.Fatalf("Unable to get VMClarity server address: %v", err)
		}
		assetScanID, err := createAssetScan(context.TODO(), server, assetID)
		if err != nil {
			logger.Fatalf("Failed to create asset scan: %v", err)
		}
		fmt.Println(assetScanID)
	},
}

func init() {
	rootCmd.AddCommand(assetScanCreateCmd)
	assetScanCreateCmd.Flags().String("server", "", "VMClarity server to create asset to, for example: http://localhost:9999/api")
	assetScanCreateCmd.Flags().String("asset-id", "", "Asset ID for asset scan")
	assetScanCreateCmd.MarkFlagRequired("server")
	assetScanCreateCmd.MarkFlagRequired("asset-id")
}

func createAssetScan(ctx context.Context, server, assetID string) (string, error) {
	client, err := backendclient.Create(server)
	if err != nil {
		return "", fmt.Errorf("failed to create VMClarity API client: %w", err)
	}

	asset, err := client.GetAsset(ctx, assetID, models.GetAssetsAssetIDParams{})
	if err != nil {
		return "", fmt.Errorf("failed to get asset %s: %w", assetID, err)
	}
	assetScanData := createEmptyAssetScanForAsset(asset)

	assetScan, err := client.PostAssetScan(ctx, assetScanData)
	if err != nil {
		var conErr backendclient.AssetScanConflictError
		if errors.As(err, &conErr) {
			assetScanID := *conErr.ConflictingAssetScan.Id
			logger.WithField("AssetScanID", assetScanID).Debug("AssetScan already exist.")
			return *conErr.ConflictingAssetScan.Id, nil
		}
		return "", fmt.Errorf("failed to post AssetScan to backend API: %w", err)
	}

	return *assetScan.Id, nil
}

func createEmptyAssetScanForAsset(asset models.Asset) models.AssetScan {
	return models.AssetScan{
		Asset: &models.AssetRelationship{
			AssetInfo: asset.AssetInfo,
			FirstSeen: asset.FirstSeen,
			Id:        *asset.Id,
		},
		Status: &models.AssetScanStatus{
			General: &models.AssetScanState{
				State: utils.PointerTo(models.AssetScanStateStateReadyToScan),
			},
		},
	}
}
