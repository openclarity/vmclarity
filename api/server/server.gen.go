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
	GetTargetsTargetID(ctx echo.Context, targetID TargetID) error
	// Update target.
	// (PUT /targets/{targetID})
	PutTargetsTargetID(ctx echo.Context, targetID TargetID) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
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
