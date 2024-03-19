package client

import (
	"context"
	"encoding/json"
	internal "github.com/openclarity/vmclarity/scanner/client/internal/client"
	"github.com/openclarity/vmclarity/scanner/types"
	"io"
	"net/http"
)

// RequestEditorFn  is the function signature for the RequestEditor callback function
type RequestEditorFn = internal.RequestEditorFn

type Client struct {
	api *internal.Client
}

func NewClient(server string) (*Client, error) {
	c, err := internal.NewClient(server)
	if err != nil {
		return nil, err
	}
	return &Client{c}, nil
}

func (c *Client) GetFindings(ctx context.Context, params *types.GetFindingsParams, reqEditors ...RequestEditorFn) (*types.ScanFindings, error) {
	return respTo[*types.ScanFindings](c.api.GetFindings(ctx, params, reqEditors...))
}

func (c *Client) GetScanFindingsForScan(ctx context.Context, scanID string, reqEditors ...RequestEditorFn) (*types.ScanFindings, error) {
	return respTo[*types.ScanFindings](c.api.GetScanFindingsForScan(ctx, scanID, reqEditors...))
}

func (c *Client) IsAlive(ctx context.Context, reqEditors ...RequestEditorFn) (string, error) {
	return respTo[string](c.api.IsAlive(ctx, reqEditors...))
}

func (c *Client) IsReady(ctx context.Context, reqEditors ...RequestEditorFn) (string, error) {
	return respTo[string](c.api.IsReady(ctx, reqEditors...))
}

func (c *Client) GetScan(ctx context.Context, scanID string, reqEditors ...RequestEditorFn) (*types.Scan, error) {
	return respTo[*types.Scan](c.api.GetScan(ctx, scanID, reqEditors...))
}

func (c *Client) SubmitScanEvent(ctx context.Context, scanID string, event types.ScanEvent, reqEditors ...RequestEditorFn) (*types.Scan, error) {
	return respTo[*types.Scan](c.api.SubmitScanEvent(ctx, scanID, types.SubmitScanEventJSONRequestBody{
		EventInfo: event.EventInfo,
	}, reqEditors...))
}

func (c *Client) MarkScanAborted(ctx context.Context, scanID string, reqEditors ...RequestEditorFn) (*types.Scan, error) {
	return respTo[*types.Scan](c.api.MarkScanAborted(ctx, scanID, reqEditors...))
}

func (c *Client) GetScans(ctx context.Context, params types.GetScansParams, reqEditors ...RequestEditorFn) (*types.Scans, error) {
	return respTo[*types.Scans](c.api.GetScans(ctx, &params, reqEditors...))
}

func (c *Client) CreateScan(ctx context.Context, scan types.Scan, body io.Reader, reqEditors ...RequestEditorFn) (*types.Scan, error) {
	return respTo[*types.Scan](c.api.CreateScan(ctx, scan, reqEditors...))
}

func respTo[T any](r *http.Response, err error) (T, error) {
	var target T

	// check original error
	if err != nil {
		return target, err
	}

	// decode
	defer r.Body.Close()
	err = json.NewDecoder(r.Body).Decode(&target)
	if err != nil {
		return target, err
	}

	return target, nil
}
