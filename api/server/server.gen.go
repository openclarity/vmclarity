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
	"9mChU7IUnsETVBC4+rD6/wAAAP//Jh37U89uAAA=",
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
