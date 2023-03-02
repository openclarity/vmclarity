package backend_client

import (
	"context"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/openclarity/vmclarity/api/client"
	"github.com/openclarity/vmclarity/api/models"
	runtimeScanUtils "github.com/openclarity/vmclarity/runtime_scan/pkg/utils"
	"github.com/openclarity/vmclarity/shared/pkg/utils"
)

type BackendClient struct {
	apiClient client.ClientWithResponsesInterface
}

func Create(serverAddress string) (*BackendClient, error) {
	apiClient, err := client.NewClientWithResponses(serverAddress)
	if err != nil {
		return nil, fmt.Errorf("unable to create VMClarity API client. serverAddress=%v: %w", serverAddress, err)
	}

	return &BackendClient{
		apiClient: apiClient,
	}, nil
}

func (b *BackendClient) GetScanResult(scanResultID string) (models.TargetScanResult, error) {
	newGetExistingError := func(err error) error {
		return fmt.Errorf("failed to get existing scan result %v: %w", scanResultID, err)
	}

	var scanResults models.TargetScanResult
	resp, err := b.apiClient.GetScanResultsScanResultIDWithResponse(context.TODO(), scanResultID, &models.GetScanResultsScanResultIDParams{})
	if err != nil {
		return scanResults, newGetExistingError(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		if resp.JSON200 == nil {
			return scanResults, newGetExistingError(fmt.Errorf("empty body"))
		}
		return *resp.JSON200, nil
	case http.StatusNotFound:
		if resp.JSON404 == nil {
			return scanResults, newGetExistingError(fmt.Errorf("empty body on not found"))
		}
		if resp.JSON404 != nil && resp.JSON404.Message != nil {
			return scanResults, newGetExistingError(fmt.Errorf("not found: %v", *resp.JSON404.Message))
		}
		return scanResults, newGetExistingError(fmt.Errorf("not found"))
	default:
		if resp.JSONDefault != nil && resp.JSONDefault.Message != nil {
			return scanResults, newGetExistingError(fmt.Errorf("status code=%v: %v", resp.StatusCode(), *resp.JSONDefault.Message))
		}
		return scanResults, newGetExistingError(fmt.Errorf("status code=%v", resp.StatusCode()))
	}
}

func (b *BackendClient) PatchScanResult(scanResult models.TargetScanResult, scanResultID string) error {
	newUpdateScanResultError := func(err error) error {
		return fmt.Errorf("failed to update scan result %v: %w", scanResultID, err)
	}

	resp, err := b.apiClient.PatchScanResultsScanResultIDWithResponse(context.TODO(), scanResultID, scanResult)
	if err != nil {
		return newUpdateScanResultError(err)
	}

	switch resp.StatusCode() {
	case http.StatusOK:
		if resp.JSON200 == nil {
			return newUpdateScanResultError(fmt.Errorf("empty body"))
		}
		return nil
	case http.StatusNotFound:
		if resp.JSON404 == nil {
			return newUpdateScanResultError(fmt.Errorf("empty body on not found"))
		}
		if resp.JSON404 != nil && resp.JSON404.Message != nil {
			return newUpdateScanResultError(fmt.Errorf("not found: %v", *resp.JSON404.Message))
		}
		return newUpdateScanResultError(fmt.Errorf("not found"))
	default:
		if resp.JSONDefault != nil && resp.JSONDefault.Message != nil {
			return newUpdateScanResultError(fmt.Errorf("status code=%v: %v", resp.StatusCode(), *resp.JSONDefault.Message))
		}
		return newUpdateScanResultError(fmt.Errorf("status code=%v", resp.StatusCode()))
	}
}

func (b *BackendClient) PostScan(ctx context.Context, scan models.Scan) (string, error) {
	resp, err := b.apiClient.PostScansWithResponse(ctx, scan)
	if err != nil {
		return "", fmt.Errorf("failed to post a scan: %v", err)
	}
	switch resp.StatusCode() {
	case http.StatusCreated:
		if resp.JSON201 == nil {
			return "", fmt.Errorf("failed to create a scan: empty body")
		}
		if resp.JSON201.Id == nil {
			return "", fmt.Errorf("scan id is nil")
		}
		return *resp.JSON201.Id, nil
	case http.StatusConflict:
		if resp.JSON409 == nil {
			return "", fmt.Errorf("failed to create a scan: empty body on conflict")
		}
		if resp.JSON409.Scan.Id == nil {
			return "", fmt.Errorf("scan id on conflict is nil")
		}
		return *resp.JSON409.Scan.Id, nil
	default:
		if resp.JSONDefault != nil && resp.JSONDefault.Message != nil {
			return "", fmt.Errorf("failed to post scan. status code=%v: %v", resp.StatusCode(), resp.JSONDefault.Message)
		}
		return "", fmt.Errorf("failed to post scan. status code=%v", resp.StatusCode())
	}
}

func (b *BackendClient) PostScanResult(ctx context.Context, scanResult models.TargetScanResult) (string, error) {
	resp, err := b.apiClient.PostScanResultsWithResponse(ctx, scanResult)
	if err != nil {
		return "", fmt.Errorf("failed to post scan status: %v", err)
	}
	switch resp.StatusCode() {
	case http.StatusCreated:
		if resp.JSON201 == nil {
			return "", fmt.Errorf("failed to create a scan status, empty body")
		}
		if resp.JSON201.Id == nil {
			return "", fmt.Errorf("failed to create a scan status, missing id")
		}
		return *resp.JSON201.Id, nil
	case http.StatusConflict:
		if resp.JSON409 == nil {
			return "", fmt.Errorf("failed to create a scan status, empty body on conflict")
		}
		if resp.JSON409.TargetScanResult.Id == nil {
			return "", fmt.Errorf("failed to create a scan status, missing id")
		}
		log.Infof("Scan results already exist with id %v.", *resp.JSON409.TargetScanResult.Id)
		return *resp.JSON409.TargetScanResult.Id, nil
	default:
		if resp.JSONDefault != nil && resp.JSONDefault.Message != nil {
			return "", fmt.Errorf("failed to create a scan status. status code=%v: %v", resp.StatusCode(), resp.JSONDefault.Message)
		}
		return "", fmt.Errorf("failed to create a scan status. status code=%v", resp.StatusCode())
	}
}

func (b *BackendClient) PatchScan(ctx context.Context, scanID models.ScanID, scan *models.Scan) error {
	resp, err := b.apiClient.PatchScansScanIDWithResponse(ctx, scanID, *scan)
	if err != nil {
		return fmt.Errorf("failed to patch a scan: %v", err)
	}
	switch resp.StatusCode() {
	case http.StatusOK:
		if resp.JSON200 == nil {
			return fmt.Errorf("failed to patch a scan: empty body")
		}
		return nil
	default:
		if resp.JSONDefault != nil && resp.JSONDefault.Message != nil {
			return fmt.Errorf("failed to patch scan. status code=%v: %v", resp.StatusCode(), resp.JSONDefault.Message)
		}
		return fmt.Errorf("failed to patch scan. status code=%v", resp.StatusCode())
	}
}

func (b *BackendClient) GetTargetScanResultSummary(ctx context.Context, scanResultID string) (*models.TargetScanResultSummary, error) {
	params := &models.GetScanResultsScanResultIDParams{
		Select: runtimeScanUtils.StringPtr("summary"),
	}
	resp, err := b.apiClient.GetScanResultsScanResultIDWithResponse(ctx, scanResultID, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get a target scan summary: %v", err)
	}
	switch resp.StatusCode() {
	case http.StatusOK:
		if resp.JSON200 == nil {
			return nil, fmt.Errorf("failed to get a target scan summary: empty body")
		}
		if resp.JSON200.Summary == nil {
			return nil, fmt.Errorf("failed to get a target scan summary: empty summary in body")
		}
		return resp.JSON200.Summary, nil
	default:
		if resp.JSONDefault != nil && resp.JSONDefault.Message != nil {
			return nil, fmt.Errorf("failed to get a target scan summary. summary code=%v: %v", resp.StatusCode(), resp.JSONDefault.Message)
		}
		return nil, fmt.Errorf("failed to get a target scan summary. summary code=%v", resp.StatusCode())
	}
}

func (b *BackendClient) GetTargetScanStatus(ctx context.Context, scanResultID string) (*models.TargetScanStatus, error) {
	params := &models.GetScanResultsScanResultIDParams{
		Select: utils.StringPtr("status"),
	}
	resp, err := b.apiClient.GetScanResultsScanResultIDWithResponse(ctx, scanResultID, params)
	if err != nil {
		return nil, fmt.Errorf("failed to get a target scan status: %v", err)
	}
	switch resp.StatusCode() {
	case http.StatusOK:
		if resp.JSON200 == nil {
			return nil, fmt.Errorf("failed to get a target scan status: empty body")
		}
		if resp.JSON200.Status == nil {
			return nil, fmt.Errorf("failed to get a target scan status: empty status in body")
		}
		return resp.JSON200.Status, nil
	default:
		if resp.JSONDefault != nil && resp.JSONDefault.Message != nil {
			return nil, fmt.Errorf("failed to get a target scan status. status code=%v: %v", resp.StatusCode(), resp.JSONDefault.Message)
		}
		return nil, fmt.Errorf("failed to get a target scan status. status code=%v", resp.StatusCode())
	}
}

func (b *BackendClient) PatchTargetScanStatus(ctx context.Context, scanResultID string, status *models.TargetScanStatus) error {
	scanResult := models.TargetScanResult{
		Status: status,
	}
	resp, err := b.apiClient.PatchScanResultsScanResultIDWithResponse(ctx, scanResultID, scanResult)
	if err != nil {
		return fmt.Errorf("failed to patch a scan result status: %v", err)
	}
	switch resp.StatusCode() {
	case http.StatusOK:
		if resp.JSON200 == nil {
			return fmt.Errorf("failed to update a scan result status: empty body")
		}
		return nil
	case http.StatusNotFound:
		if resp.JSON404 == nil {
			return fmt.Errorf("failed to update a scan result status: empty body on not found")
		}
		if resp.JSON404 != nil && resp.JSON404.Message != nil {
			return fmt.Errorf("failed to update scan result status, not found: %v", resp.JSON404.Message)
		}
		return fmt.Errorf("failed to update scan result status, not found")
	default:
		if resp.JSONDefault != nil && resp.JSONDefault.Message != nil {
			return fmt.Errorf("failed to update scan result status. status code=%v: %v", resp.StatusCode(), resp.JSONDefault.Message)
		}
		return fmt.Errorf("failed to update scan result status. status code=%v", resp.StatusCode())
	}
}

func (b *BackendClient) GetScan(ctx context.Context, scanID string) (*models.Scan, error) {
	resp, err := b.apiClient.GetScansScanIDWithResponse(ctx, scanID)
	if err != nil {
		return nil, fmt.Errorf("failed to get a scan: %v", err)
	}
	switch resp.StatusCode() {
	case http.StatusOK:
		if resp.JSON200 == nil {
			return nil, fmt.Errorf("failed to get a scan: empty body")
		}
		return resp.JSON200, nil
	default:
		if resp.JSONDefault != nil && resp.JSONDefault.Message != nil {
			return nil, fmt.Errorf("failed to get a scan status. status code=%v: %v", resp.StatusCode(), resp.JSONDefault.Message)
		}
		return nil, fmt.Errorf("failed to get a scan status. status code=%v", resp.StatusCode())
	}
}

func (b *BackendClient) GetScanConfigs(ctx context.Context) (*models.ScanConfigs, error) {
	resp, err := b.apiClient.GetScanConfigsWithResponse(ctx, &models.GetScanConfigsParams{})
	if err != nil {
		return nil, fmt.Errorf("failed to get scan configs: %v", err)
	}
	switch resp.StatusCode() {
	case http.StatusOK:
		if resp.JSON200 == nil {
			return nil, fmt.Errorf("no scan configs: empty body")
		}
		return resp.JSON200, nil
	default:
		if resp.JSONDefault != nil && resp.JSONDefault.Message != nil {
			return nil, fmt.Errorf("failed to get scan configs. status code=%v: %s", resp.StatusCode(), *resp.JSONDefault.Message)
		}
		return nil, fmt.Errorf("failed to get scan configs. status code=%v", resp.StatusCode())
	}
}

func (b *BackendClient) GetScans(ctx context.Context) (*models.Scans, error) {
	resp, err := b.apiClient.GetScansWithResponse(ctx, &models.GetScansParams{})
	if err != nil {
		return nil, fmt.Errorf("failed to get scans: %v", err)
	}
	switch resp.StatusCode() {
	case http.StatusOK:
		if resp.JSON200 == nil {
			return nil, fmt.Errorf("no scans: empty body")
		}
	default:
		if resp.JSONDefault != nil && resp.JSONDefault.Message != nil {
			return nil, fmt.Errorf("failed to get scans. status code=%v: %s", resp.StatusCode(), *resp.JSONDefault.Message)
		}
		return nil, fmt.Errorf("failed to get scans. status code=%v", resp.StatusCode())
	}
	return resp.JSON200, nil
}

func (b *BackendClient) PostTarget(ctx context.Context, target models.Target) (string, error) {
	resp, err := b.apiClient.PostTargetsWithResponse(ctx, target)
	if err != nil {
		return "", fmt.Errorf("failed to post target: %v", err)
	}
	switch resp.StatusCode() {
	case http.StatusCreated:
		if resp.JSON201 == nil {
			return "", fmt.Errorf("failed to create a target: empty body")
		}
		return *resp.JSON201.Id, nil
	case http.StatusConflict:
		if resp.JSON409 == nil {
			return "", fmt.Errorf("failed to create a target: empty body on conflict")
		}
		return *resp.JSON409.Target.Id, nil
	default:
		if resp.JSONDefault != nil && resp.JSONDefault.Message != nil {
			return "", fmt.Errorf("failed to post target. status code=%v: %v", resp.StatusCode(), resp.JSONDefault.Message)
		}
		return "", fmt.Errorf("failed to post target. status code=%v", resp.StatusCode())
	}
}
