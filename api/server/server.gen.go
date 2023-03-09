// Package server provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.3 DO NOT EDIT.
package server

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	. "github.com/openclarity/vmclarity/api/models"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Get all available scopes
	// (GET /discovery/scopes)
	GetDiscoveryScopes(ctx echo.Context, params GetDiscoveryScopesParams) error
	// Set all available scopes
	// (PUT /discovery/scopes)
	PutDiscoveryScopes(ctx echo.Context) error
	// Get all scan configs.
	// (GET /scanConfigs)
	GetScanConfigs(ctx echo.Context, params GetScanConfigsParams) error
	// Create a scan config
	// (POST /scanConfigs)
	PostScanConfigs(ctx echo.Context) error
	// Delete a scan config.
	// (DELETE /scanConfigs/{scanConfigID})
	DeleteScanConfigsScanConfigID(ctx echo.Context, scanConfigID ScanConfigID) error
	// Get the details for a scan config.
	// (GET /scanConfigs/{scanConfigID})
	GetScanConfigsScanConfigID(ctx echo.Context, scanConfigID ScanConfigID, params GetScanConfigsScanConfigIDParams) error
	// Patch a scan config.
	// (PATCH /scanConfigs/{scanConfigID})
	PatchScanConfigsScanConfigID(ctx echo.Context, scanConfigID ScanConfigID) error
	// Update a scan config.
	// (PUT /scanConfigs/{scanConfigID})
	PutScanConfigsScanConfigID(ctx echo.Context, scanConfigID ScanConfigID) error
	// Get scan results according to the given filters
	// (GET /scanResults)
	GetScanResults(ctx echo.Context, params GetScanResultsParams) error
	// Create a scan result for a specific target for a specific scan
	// (POST /scanResults)
	PostScanResults(ctx echo.Context) error
	// Get a scan result.
	// (GET /scanResults/{scanResultID})
	GetScanResultsScanResultID(ctx echo.Context, scanResultID ScanResultID, params GetScanResultsScanResultIDParams) error
	// Patch a scan result
	// (PATCH /scanResults/{scanResultID})
	PatchScanResultsScanResultID(ctx echo.Context, scanResultID ScanResultID) error
	// Update a scan result.
	// (PUT /scanResults/{scanResultID})
	PutScanResultsScanResultID(ctx echo.Context, scanResultID ScanResultID) error
	// Get all scans. Each scan contaians details about a multi-target scheduled scan.
	// (GET /scans)
	GetScans(ctx echo.Context, params GetScansParams) error
	// Create a multi-target scheduled scan
	// (POST /scans)
	PostScans(ctx echo.Context) error
	// Delete a scan.
	// (DELETE /scans/{scanID})
	DeleteScansScanID(ctx echo.Context, scanID ScanID) error
	// Get the details for a given multi-target scheduled scan.
	// (GET /scans/{scanID})
	GetScansScanID(ctx echo.Context, scanID ScanID) error
	// Patch a scan.
	// (PATCH /scans/{scanID})
	PatchScansScanID(ctx echo.Context, scanID ScanID) error
	// Update a scan.
	// (PUT /scans/{scanID})
	PutScansScanID(ctx echo.Context, scanID ScanID) error
	// Get targets
	// (GET /targets)
	GetTargets(ctx echo.Context, params GetTargetsParams) error
	// Create target
	// (POST /targets)
	PostTargets(ctx echo.Context) error
	// Delete target.
	// (DELETE /targets/{targetID})
	DeleteTargetsTargetID(ctx echo.Context, targetID TargetID) error
	// Get target.
	// (GET /targets/{targetID})
	GetTargetsTargetID(ctx echo.Context, targetID TargetID, params GetTargetsTargetIDParams) error
	// Update target.
	// (PUT /targets/{targetID})
	PutTargetsTargetID(ctx echo.Context, targetID TargetID) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetDiscoveryScopes converts echo context to params.
func (w *ServerInterfaceWrapper) GetDiscoveryScopes(ctx echo.Context) error {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetDiscoveryScopesParams
	// ------------- Optional query parameter "$filter" -------------

	err = runtime.BindQueryParameter("form", true, false, "$filter", ctx.QueryParams(), &params.Filter)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter $filter: %s", err))
	}

	// ------------- Optional query parameter "$select" -------------

	err = runtime.BindQueryParameter("form", true, false, "$select", ctx.QueryParams(), &params.Select)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter $select: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetDiscoveryScopes(ctx, params)
	return err
}

// PutDiscoveryScopes converts echo context to params.
func (w *ServerInterfaceWrapper) PutDiscoveryScopes(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PutDiscoveryScopes(ctx)
	return err
}

// GetScanConfigs converts echo context to params.
func (w *ServerInterfaceWrapper) GetScanConfigs(ctx echo.Context) error {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetScanConfigsParams
	// ------------- Optional query parameter "$filter" -------------

	err = runtime.BindQueryParameter("form", true, false, "$filter", ctx.QueryParams(), &params.Filter)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter $filter: %s", err))
	}

	// ------------- Optional query parameter "$select" -------------

	err = runtime.BindQueryParameter("form", true, false, "$select", ctx.QueryParams(), &params.Select)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter $select: %s", err))
	}

	// ------------- Optional query parameter "$count" -------------

	err = runtime.BindQueryParameter("form", true, false, "$count", ctx.QueryParams(), &params.Count)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter $count: %s", err))
	}

	// ------------- Optional query parameter "$top" -------------

	err = runtime.BindQueryParameter("form", true, false, "$top", ctx.QueryParams(), &params.Top)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter $top: %s", err))
	}

	// ------------- Optional query parameter "$skip" -------------

	err = runtime.BindQueryParameter("form", true, false, "$skip", ctx.QueryParams(), &params.Skip)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter $skip: %s", err))
	}

	// ------------- Optional query parameter "$expand" -------------

	err = runtime.BindQueryParameter("form", true, false, "$expand", ctx.QueryParams(), &params.Expand)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter $expand: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetScanConfigs(ctx, params)
	return err
}

// PostScanConfigs converts echo context to params.
func (w *ServerInterfaceWrapper) PostScanConfigs(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PostScanConfigs(ctx)
	return err
}

// DeleteScanConfigsScanConfigID converts echo context to params.
func (w *ServerInterfaceWrapper) DeleteScanConfigsScanConfigID(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "scanConfigID" -------------
	var scanConfigID ScanConfigID

	err = runtime.BindStyledParameterWithLocation("simple", false, "scanConfigID", runtime.ParamLocationPath, ctx.Param("scanConfigID"), &scanConfigID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter scanConfigID: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.DeleteScanConfigsScanConfigID(ctx, scanConfigID)
	return err
}

// GetScanConfigsScanConfigID converts echo context to params.
func (w *ServerInterfaceWrapper) GetScanConfigsScanConfigID(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "scanConfigID" -------------
	var scanConfigID ScanConfigID

	err = runtime.BindStyledParameterWithLocation("simple", false, "scanConfigID", runtime.ParamLocationPath, ctx.Param("scanConfigID"), &scanConfigID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter scanConfigID: %s", err))
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params GetScanConfigsScanConfigIDParams
	// ------------- Optional query parameter "$select" -------------

	err = runtime.BindQueryParameter("form", true, false, "$select", ctx.QueryParams(), &params.Select)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter $select: %s", err))
	}

	// ------------- Optional query parameter "$expand" -------------

	err = runtime.BindQueryParameter("form", true, false, "$expand", ctx.QueryParams(), &params.Expand)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter $expand: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetScanConfigsScanConfigID(ctx, scanConfigID, params)
	return err
}

// PatchScanConfigsScanConfigID converts echo context to params.
func (w *ServerInterfaceWrapper) PatchScanConfigsScanConfigID(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "scanConfigID" -------------
	var scanConfigID ScanConfigID

	err = runtime.BindStyledParameterWithLocation("simple", false, "scanConfigID", runtime.ParamLocationPath, ctx.Param("scanConfigID"), &scanConfigID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter scanConfigID: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PatchScanConfigsScanConfigID(ctx, scanConfigID)
	return err
}

// PutScanConfigsScanConfigID converts echo context to params.
func (w *ServerInterfaceWrapper) PutScanConfigsScanConfigID(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "scanConfigID" -------------
	var scanConfigID ScanConfigID

	err = runtime.BindStyledParameterWithLocation("simple", false, "scanConfigID", runtime.ParamLocationPath, ctx.Param("scanConfigID"), &scanConfigID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter scanConfigID: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PutScanConfigsScanConfigID(ctx, scanConfigID)
	return err
}

// GetScanResults converts echo context to params.
func (w *ServerInterfaceWrapper) GetScanResults(ctx echo.Context) error {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetScanResultsParams
	// ------------- Optional query parameter "$filter" -------------

	err = runtime.BindQueryParameter("form", true, false, "$filter", ctx.QueryParams(), &params.Filter)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter $filter: %s", err))
	}

	// ------------- Optional query parameter "$select" -------------

	err = runtime.BindQueryParameter("form", true, false, "$select", ctx.QueryParams(), &params.Select)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter $select: %s", err))
	}

	// ------------- Optional query parameter "$count" -------------

	err = runtime.BindQueryParameter("form", true, false, "$count", ctx.QueryParams(), &params.Count)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter $count: %s", err))
	}

	// ------------- Optional query parameter "$top" -------------

	err = runtime.BindQueryParameter("form", true, false, "$top", ctx.QueryParams(), &params.Top)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter $top: %s", err))
	}

	// ------------- Optional query parameter "$skip" -------------

	err = runtime.BindQueryParameter("form", true, false, "$skip", ctx.QueryParams(), &params.Skip)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter $skip: %s", err))
	}

	// ------------- Optional query parameter "$expand" -------------

	err = runtime.BindQueryParameter("form", true, false, "$expand", ctx.QueryParams(), &params.Expand)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter $expand: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetScanResults(ctx, params)
	return err
}

// PostScanResults converts echo context to params.
func (w *ServerInterfaceWrapper) PostScanResults(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PostScanResults(ctx)
	return err
}

// GetScanResultsScanResultID converts echo context to params.
func (w *ServerInterfaceWrapper) GetScanResultsScanResultID(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "scanResultID" -------------
	var scanResultID ScanResultID

	err = runtime.BindStyledParameterWithLocation("simple", false, "scanResultID", runtime.ParamLocationPath, ctx.Param("scanResultID"), &scanResultID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter scanResultID: %s", err))
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params GetScanResultsScanResultIDParams
	// ------------- Optional query parameter "$select" -------------

	err = runtime.BindQueryParameter("form", true, false, "$select", ctx.QueryParams(), &params.Select)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter $select: %s", err))
	}

	// ------------- Optional query parameter "$expand" -------------

	err = runtime.BindQueryParameter("form", true, false, "$expand", ctx.QueryParams(), &params.Expand)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter $expand: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetScanResultsScanResultID(ctx, scanResultID, params)
	return err
}

// PatchScanResultsScanResultID converts echo context to params.
func (w *ServerInterfaceWrapper) PatchScanResultsScanResultID(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "scanResultID" -------------
	var scanResultID ScanResultID

	err = runtime.BindStyledParameterWithLocation("simple", false, "scanResultID", runtime.ParamLocationPath, ctx.Param("scanResultID"), &scanResultID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter scanResultID: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PatchScanResultsScanResultID(ctx, scanResultID)
	return err
}

// PutScanResultsScanResultID converts echo context to params.
func (w *ServerInterfaceWrapper) PutScanResultsScanResultID(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "scanResultID" -------------
	var scanResultID ScanResultID

	err = runtime.BindStyledParameterWithLocation("simple", false, "scanResultID", runtime.ParamLocationPath, ctx.Param("scanResultID"), &scanResultID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter scanResultID: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PutScanResultsScanResultID(ctx, scanResultID)
	return err
}

// GetScans converts echo context to params.
func (w *ServerInterfaceWrapper) GetScans(ctx echo.Context) error {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetScansParams
	// ------------- Optional query parameter "$filter" -------------

	err = runtime.BindQueryParameter("form", true, false, "$filter", ctx.QueryParams(), &params.Filter)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter $filter: %s", err))
	}

	// ------------- Optional query parameter "page" -------------

	err = runtime.BindQueryParameter("form", true, false, "page", ctx.QueryParams(), &params.Page)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter page: %s", err))
	}

	// ------------- Optional query parameter "pageSize" -------------

	err = runtime.BindQueryParameter("form", true, false, "pageSize", ctx.QueryParams(), &params.PageSize)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter pageSize: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetScans(ctx, params)
	return err
}

// PostScans converts echo context to params.
func (w *ServerInterfaceWrapper) PostScans(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PostScans(ctx)
	return err
}

// DeleteScansScanID converts echo context to params.
func (w *ServerInterfaceWrapper) DeleteScansScanID(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "scanID" -------------
	var scanID ScanID

	err = runtime.BindStyledParameterWithLocation("simple", false, "scanID", runtime.ParamLocationPath, ctx.Param("scanID"), &scanID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter scanID: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.DeleteScansScanID(ctx, scanID)
	return err
}

// GetScansScanID converts echo context to params.
func (w *ServerInterfaceWrapper) GetScansScanID(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "scanID" -------------
	var scanID ScanID

	err = runtime.BindStyledParameterWithLocation("simple", false, "scanID", runtime.ParamLocationPath, ctx.Param("scanID"), &scanID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter scanID: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetScansScanID(ctx, scanID)
	return err
}

// PatchScansScanID converts echo context to params.
func (w *ServerInterfaceWrapper) PatchScansScanID(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "scanID" -------------
	var scanID ScanID

	err = runtime.BindStyledParameterWithLocation("simple", false, "scanID", runtime.ParamLocationPath, ctx.Param("scanID"), &scanID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter scanID: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PatchScansScanID(ctx, scanID)
	return err
}

// PutScansScanID converts echo context to params.
func (w *ServerInterfaceWrapper) PutScansScanID(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "scanID" -------------
	var scanID ScanID

	err = runtime.BindStyledParameterWithLocation("simple", false, "scanID", runtime.ParamLocationPath, ctx.Param("scanID"), &scanID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter scanID: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PutScansScanID(ctx, scanID)
	return err
}

// GetTargets converts echo context to params.
func (w *ServerInterfaceWrapper) GetTargets(ctx echo.Context) error {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetTargetsParams
	// ------------- Optional query parameter "$filter" -------------

	err = runtime.BindQueryParameter("form", true, false, "$filter", ctx.QueryParams(), &params.Filter)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter $filter: %s", err))
	}

	// ------------- Optional query parameter "$select" -------------

	err = runtime.BindQueryParameter("form", true, false, "$select", ctx.QueryParams(), &params.Select)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter $select: %s", err))
	}

	// ------------- Optional query parameter "$count" -------------

	err = runtime.BindQueryParameter("form", true, false, "$count", ctx.QueryParams(), &params.Count)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter $count: %s", err))
	}

	// ------------- Optional query parameter "$top" -------------

	err = runtime.BindQueryParameter("form", true, false, "$top", ctx.QueryParams(), &params.Top)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter $top: %s", err))
	}

	// ------------- Optional query parameter "$skip" -------------

	err = runtime.BindQueryParameter("form", true, false, "$skip", ctx.QueryParams(), &params.Skip)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter $skip: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetTargets(ctx, params)
	return err
}

// PostTargets converts echo context to params.
func (w *ServerInterfaceWrapper) PostTargets(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PostTargets(ctx)
	return err
}

// DeleteTargetsTargetID converts echo context to params.
func (w *ServerInterfaceWrapper) DeleteTargetsTargetID(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "targetID" -------------
	var targetID TargetID

	err = runtime.BindStyledParameterWithLocation("simple", false, "targetID", runtime.ParamLocationPath, ctx.Param("targetID"), &targetID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter targetID: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.DeleteTargetsTargetID(ctx, targetID)
	return err
}

// GetTargetsTargetID converts echo context to params.
func (w *ServerInterfaceWrapper) GetTargetsTargetID(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "targetID" -------------
	var targetID TargetID

	err = runtime.BindStyledParameterWithLocation("simple", false, "targetID", runtime.ParamLocationPath, ctx.Param("targetID"), &targetID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter targetID: %s", err))
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params GetTargetsTargetIDParams
	// ------------- Optional query parameter "$select" -------------

	err = runtime.BindQueryParameter("form", true, false, "$select", ctx.QueryParams(), &params.Select)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter $select: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetTargetsTargetID(ctx, targetID, params)
	return err
}

// PutTargetsTargetID converts echo context to params.
func (w *ServerInterfaceWrapper) PutTargetsTargetID(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "targetID" -------------
	var targetID TargetID

	err = runtime.BindStyledParameterWithLocation("simple", false, "targetID", runtime.ParamLocationPath, ctx.Param("targetID"), &targetID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter targetID: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PutTargetsTargetID(ctx, targetID)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/discovery/scopes", wrapper.GetDiscoveryScopes)
	router.PUT(baseURL+"/discovery/scopes", wrapper.PutDiscoveryScopes)
	router.GET(baseURL+"/scanConfigs", wrapper.GetScanConfigs)
	router.POST(baseURL+"/scanConfigs", wrapper.PostScanConfigs)
	router.DELETE(baseURL+"/scanConfigs/:scanConfigID", wrapper.DeleteScanConfigsScanConfigID)
	router.GET(baseURL+"/scanConfigs/:scanConfigID", wrapper.GetScanConfigsScanConfigID)
	router.PATCH(baseURL+"/scanConfigs/:scanConfigID", wrapper.PatchScanConfigsScanConfigID)
	router.PUT(baseURL+"/scanConfigs/:scanConfigID", wrapper.PutScanConfigsScanConfigID)
	router.GET(baseURL+"/scanResults", wrapper.GetScanResults)
	router.POST(baseURL+"/scanResults", wrapper.PostScanResults)
	router.GET(baseURL+"/scanResults/:scanResultID", wrapper.GetScanResultsScanResultID)
	router.PATCH(baseURL+"/scanResults/:scanResultID", wrapper.PatchScanResultsScanResultID)
	router.PUT(baseURL+"/scanResults/:scanResultID", wrapper.PutScanResultsScanResultID)
	router.GET(baseURL+"/scans", wrapper.GetScans)
	router.POST(baseURL+"/scans", wrapper.PostScans)
	router.DELETE(baseURL+"/scans/:scanID", wrapper.DeleteScansScanID)
	router.GET(baseURL+"/scans/:scanID", wrapper.GetScansScanID)
	router.PATCH(baseURL+"/scans/:scanID", wrapper.PatchScansScanID)
	router.PUT(baseURL+"/scans/:scanID", wrapper.PutScansScanID)
	router.GET(baseURL+"/targets", wrapper.GetTargets)
	router.POST(baseURL+"/targets", wrapper.PostTargets)
	router.DELETE(baseURL+"/targets/:targetID", wrapper.DeleteTargetsTargetID)
	router.GET(baseURL+"/targets/:targetID", wrapper.GetTargetsTargetID)
	router.PUT(baseURL+"/targets/:targetID", wrapper.PutTargetsTargetID)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
	"U1JQ/PTp6f8DAAD//yaAtEjbewAA",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
