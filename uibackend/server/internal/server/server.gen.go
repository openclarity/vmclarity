// Package server provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.2.0 DO NOT EDIT.
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

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
	. "github.com/openclarity/vmclarity/uibackend/types"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Get a list of findings impact for the dashboard.
	// (GET /dashboard/findingsImpact)
	GetDashboardFindingsImpact(ctx echo.Context) error
	// Get a list of finding trends for all finding types.
	// (GET /dashboard/findingsTrends)
	GetDashboardFindingsTrends(ctx echo.Context, params GetDashboardFindingsTrendsParams) error
	// Get a list of riskiest assets for the dashboard.
	// (GET /dashboard/riskiestAssets)
	GetDashboardRiskiestAssets(ctx echo.Context) error
	// Get a list of riskiest regions for the dashboard.
	// (GET /dashboard/riskiestRegions)
	GetDashboardRiskiestRegions(ctx echo.Context) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetDashboardFindingsImpact converts echo context to params.
func (w *ServerInterfaceWrapper) GetDashboardFindingsImpact(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetDashboardFindingsImpact(ctx)
	return err
}

// GetDashboardFindingsTrends converts echo context to params.
func (w *ServerInterfaceWrapper) GetDashboardFindingsTrends(ctx echo.Context) error {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetDashboardFindingsTrendsParams
	// ------------- Required query parameter "startTime" -------------

	err = runtime.BindQueryParameter("form", true, true, "startTime", ctx.QueryParams(), &params.StartTime)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter startTime: %s", err))
	}

	// ------------- Required query parameter "endTime" -------------

	err = runtime.BindQueryParameter("form", true, true, "endTime", ctx.QueryParams(), &params.EndTime)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter endTime: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetDashboardFindingsTrends(ctx, params)
	return err
}

// GetDashboardRiskiestAssets converts echo context to params.
func (w *ServerInterfaceWrapper) GetDashboardRiskiestAssets(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetDashboardRiskiestAssets(ctx)
	return err
}

// GetDashboardRiskiestRegions converts echo context to params.
func (w *ServerInterfaceWrapper) GetDashboardRiskiestRegions(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetDashboardRiskiestRegions(ctx)
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

	router.GET(baseURL+"/dashboard/findingsImpact", wrapper.GetDashboardFindingsImpact)
	router.GET(baseURL+"/dashboard/findingsTrends", wrapper.GetDashboardFindingsTrends)
	router.GET(baseURL+"/dashboard/riskiestAssets", wrapper.GetDashboardRiskiestAssets)
	router.GET(baseURL+"/dashboard/riskiestRegions", wrapper.GetDashboardRiskiestRegions)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
	"aysAAA==",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
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
	res := make(map[string]func() ([]byte, error))
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
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
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
