/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/openclarity/vmclarity/api/models"
	"github.com/openclarity/vmclarity/shared/pkg/backendclient"
)

// standaloneCmd represents the standalone command
var assetCreateCmd = &cobra.Command{
	Use:   "asset-create",
	Short: "Create asset",
	Long:  `It creates asset. It's useful in the CI/CD mode without VMClarity orchestration`,
	Run: func(cmd *cobra.Command, args []string) {
		logger.Infof("Creating asset...")
		filename, err := cmd.Flags().GetString("from-json-file")
		if err != nil {
			logger.Fatalf("Unable to get asset json file name: %v", err)
		}
		server, err := cmd.Flags().GetString("server")
		if err != nil {
			logger.Fatalf("Unable to get VMClarity server address: %v", err)
		}

		assetType, err := getAssetFromJsonFile(filename)
		if err != nil {
			logger.Fatalf("Failed to get asset from json file: %v", err)
		}

		_, err = assetType.ValueByDiscriminator()
		if err != nil {
			logger.Fatalf("Failed to determine asset type: %v", err)
		}

		assetID, err := createAsset(context.TODO(), assetType, server)
		if err != nil {
			logger.Fatalf("Failed to create asset: %v", err)
		}
		fmt.Println(assetID)
	},
}

func init() {
	rootCmd.AddCommand(assetCreateCmd)

	assetCreateCmd.Flags().String("from-json-file", "", "asset json filename")
	assetCreateCmd.Flags().String("server", "", "VMClarity server to create asset to, for example: http://localhost:9999/api")
	assetCreateCmd.MarkFlagRequired("from-json-file")
	assetCreateCmd.MarkFlagRequired("server")

}

func getAssetFromJsonFile(filename string) (*models.AssetType, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// get the file size
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	// read the file
	bs := make([]byte, stat.Size())
	_, err = file.Read(bs)
	if err != nil {
		return nil, err
	}

	assetType := &models.AssetType{}
	if err := assetType.UnmarshalJSON(bs); err != nil {
		return nil, fmt.Errorf("failed to unmarshal asset into AssetType %v", err)
	}

	return assetType, nil
}

func createAsset(ctx context.Context, assetType *models.AssetType, server string) (string, error) {
	client, err := backendclient.Create(server)
	if err != nil {
		return "", fmt.Errorf("failed to create VMClarity API client: %w", err)
	}

	creationTime := time.Now()
	assetData := models.Asset{
		AssetInfo: assetType,
		LastSeen:  &creationTime,
		FirstSeen: &creationTime,
	}
	asset, err := client.PostAsset(ctx, assetData)
	if err == nil {
		return *asset.Id, nil
	}
	var conflictError backendclient.AssetConflictError
	if !errors.As(err, &conflictError) {
		// If there is an error, and it's not a conflict telling
		// us that the asset already exists, then we need to
		// keep track of it and log it as a failure to
		// complete discovery. We don't fail instantly here
		// because discovering the assets is a heavy operation,
		// so we want to give the best chance to create all the
		// assets in the DB before failing.
		return "", fmt.Errorf("failed to post asset: %v", err)
	}

	// As we got a conflict it means there is an existing asset
	// which matches the unique properties of this asset, in this
	// case we'll patch the just AssetInfo and FirstSeen instead.
	assetData.FirstSeen = nil
	err = client.PatchAsset(ctx, assetData, *conflictError.ConflictingAsset.Id)
	if err != nil {
		return "", fmt.Errorf("failed to patch asset: %v", err)
	}

	return *conflictError.ConflictingAsset.Id, nil
}
