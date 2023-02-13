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
	// Discover available scopes
	// (GET /discovery/scopes)
	GetDiscoveryScopes(ctx echo.Context, params GetDiscoveryScopesParams) error
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
	GetScanConfigsScanConfigID(ctx echo.Context, scanConfigID ScanConfigID) error
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
	GetTargetsTargetID(ctx echo.Context, targetID TargetID) error
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

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetScanConfigsScanConfigID(ctx, scanConfigID)
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

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetTargetsTargetID(ctx, targetID)
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
	"AA==",
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
