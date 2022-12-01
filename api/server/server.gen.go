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
	. "github.com/vmclarity/api/models"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Get instances
	// (GET /instances)
	GetInstances(ctx echo.Context, params GetInstancesParams) error
	// Create instance
	// (POST /instances)
	PostInstances(ctx echo.Context) error
	// Delete Instance.
	// (DELETE /instances/{instanceID})
	DeleteInstancesInstanceID(ctx echo.Context, instanceID string) error
	// Get instance.
	// (GET /instances/{instanceID})
	GetInstancesInstanceID(ctx echo.Context, instanceID string) error
	// Update Application.
	// (PUT /instances/{instanceID})
	PutInstancesInstanceID(ctx echo.Context, instanceID string) error
	// Get scan results for a specified instance
	// (GET /instances/{instanceID}/scanresults)
	GetInstancesInstanceIDScanresults(ctx echo.Context, instanceID string, params GetInstancesInstanceIDScanresultsParams) error
	// Create scan result for a specified instance
	// (POST /instances/{instanceID}/scanresults)
	PostInstancesInstanceIDScanresults(ctx echo.Context, instanceID string) error
	// Report a specific scan result for a specific instance
	// (GET /instances/{instanceID}/scanresults/{scanID})
	GetInstancesInstanceIDScanresultsScanID(ctx echo.Context, instanceID string, scanID string) error
	// Report a specific exploit scan result for a specific instance
	// (GET /instances/{instanceID}/scanresults/{scanID}/exploits)
	GetInstancesInstanceIDScanresultsScanIDExploits(ctx echo.Context, instanceID string, scanID string) error
	// Report a specific malware scan result for a specific instance
	// (GET /instances/{instanceID}/scanresults/{scanID}/malwares)
	GetInstancesInstanceIDScanresultsScanIDMalwares(ctx echo.Context, instanceID string, scanID string) error
	// Report a specific misconfiguration scan result for a specific instance
	// (GET /instances/{instanceID}/scanresults/{scanID}/misconfiguration)
	GetInstancesInstanceIDScanresultsScanIDMisconfiguration(ctx echo.Context, instanceID string, scanID string) error
	// Report a specific rootkit scan result for a specific instance
	// (GET /instances/{instanceID}/scanresults/{scanID}/rootkits)
	GetInstancesInstanceIDScanresultsScanIDRootkits(ctx echo.Context, instanceID string, scanID string) error
	// Report a specific sbom scan result for a specific instance
	// (GET /instances/{instanceID}/scanresults/{scanID}/sbom)
	GetInstancesInstanceIDScanresultsScanIDSbom(ctx echo.Context, instanceID string, scanID string) error
	// Report a specific secret scan result for a specific instance
	// (GET /instances/{instanceID}/scanresults/{scanID}/secrets)
	GetInstancesInstanceIDScanresultsScanIDSecrets(ctx echo.Context, instanceID string, scanID string) error
	// Report a specific vulnerabilities scan result for a specific instance
	// (GET /instances/{instanceID}/scanresults/{scanID}/vulnerabilities)
	GetInstancesInstanceIDScanresultsScanIDVulnerabilities(ctx echo.Context, instanceID string, scanID string) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetInstances converts echo context to params.
func (w *ServerInterfaceWrapper) GetInstances(ctx echo.Context) error {
	var err error

	// Parameter object where we will unmarshal all parameters from the context
	var params GetInstancesParams
	// ------------- Required query parameter "page" -------------

	err = runtime.BindQueryParameter("form", true, true, "page", ctx.QueryParams(), &params.Page)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter page: %s", err))
	}

	// ------------- Required query parameter "pageSize" -------------

	err = runtime.BindQueryParameter("form", true, true, "pageSize", ctx.QueryParams(), &params.PageSize)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter pageSize: %s", err))
	}

	// ------------- Required query parameter "sortKey" -------------

	err = runtime.BindQueryParameter("form", true, true, "sortKey", ctx.QueryParams(), &params.SortKey)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter sortKey: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetInstances(ctx, params)
	return err
}

// PostInstances converts echo context to params.
func (w *ServerInterfaceWrapper) PostInstances(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PostInstances(ctx)
	return err
}

// DeleteInstancesInstanceID converts echo context to params.
func (w *ServerInterfaceWrapper) DeleteInstancesInstanceID(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "instanceID" -------------
	var instanceID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "instanceID", runtime.ParamLocationPath, ctx.Param("instanceID"), &instanceID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter instanceID: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.DeleteInstancesInstanceID(ctx, instanceID)
	return err
}

// GetInstancesInstanceID converts echo context to params.
func (w *ServerInterfaceWrapper) GetInstancesInstanceID(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "instanceID" -------------
	var instanceID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "instanceID", runtime.ParamLocationPath, ctx.Param("instanceID"), &instanceID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter instanceID: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetInstancesInstanceID(ctx, instanceID)
	return err
}

// PutInstancesInstanceID converts echo context to params.
func (w *ServerInterfaceWrapper) PutInstancesInstanceID(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "instanceID" -------------
	var instanceID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "instanceID", runtime.ParamLocationPath, ctx.Param("instanceID"), &instanceID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter instanceID: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PutInstancesInstanceID(ctx, instanceID)
	return err
}

// GetInstancesInstanceIDScanresults converts echo context to params.
func (w *ServerInterfaceWrapper) GetInstancesInstanceIDScanresults(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "instanceID" -------------
	var instanceID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "instanceID", runtime.ParamLocationPath, ctx.Param("instanceID"), &instanceID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter instanceID: %s", err))
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params GetInstancesInstanceIDScanresultsParams
	// ------------- Required query parameter "page" -------------

	err = runtime.BindQueryParameter("form", true, true, "page", ctx.QueryParams(), &params.Page)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter page: %s", err))
	}

	// ------------- Required query parameter "pageSize" -------------

	err = runtime.BindQueryParameter("form", true, true, "pageSize", ctx.QueryParams(), &params.PageSize)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter pageSize: %s", err))
	}

	// ------------- Required query parameter "sortKey" -------------

	err = runtime.BindQueryParameter("form", true, true, "sortKey", ctx.QueryParams(), &params.SortKey)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter sortKey: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetInstancesInstanceIDScanresults(ctx, instanceID, params)
	return err
}

// PostInstancesInstanceIDScanresults converts echo context to params.
func (w *ServerInterfaceWrapper) PostInstancesInstanceIDScanresults(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "instanceID" -------------
	var instanceID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "instanceID", runtime.ParamLocationPath, ctx.Param("instanceID"), &instanceID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter instanceID: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.PostInstancesInstanceIDScanresults(ctx, instanceID)
	return err
}

// GetInstancesInstanceIDScanresultsScanID converts echo context to params.
func (w *ServerInterfaceWrapper) GetInstancesInstanceIDScanresultsScanID(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "instanceID" -------------
	var instanceID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "instanceID", runtime.ParamLocationPath, ctx.Param("instanceID"), &instanceID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter instanceID: %s", err))
	}

	// ------------- Path parameter "scanID" -------------
	var scanID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "scanID", runtime.ParamLocationPath, ctx.Param("scanID"), &scanID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter scanID: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetInstancesInstanceIDScanresultsScanID(ctx, instanceID, scanID)
	return err
}

// GetInstancesInstanceIDScanresultsScanIDExploits converts echo context to params.
func (w *ServerInterfaceWrapper) GetInstancesInstanceIDScanresultsScanIDExploits(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "instanceID" -------------
	var instanceID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "instanceID", runtime.ParamLocationPath, ctx.Param("instanceID"), &instanceID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter instanceID: %s", err))
	}

	// ------------- Path parameter "scanID" -------------
	var scanID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "scanID", runtime.ParamLocationPath, ctx.Param("scanID"), &scanID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter scanID: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetInstancesInstanceIDScanresultsScanIDExploits(ctx, instanceID, scanID)
	return err
}

// GetInstancesInstanceIDScanresultsScanIDMalwares converts echo context to params.
func (w *ServerInterfaceWrapper) GetInstancesInstanceIDScanresultsScanIDMalwares(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "instanceID" -------------
	var instanceID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "instanceID", runtime.ParamLocationPath, ctx.Param("instanceID"), &instanceID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter instanceID: %s", err))
	}

	// ------------- Path parameter "scanID" -------------
	var scanID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "scanID", runtime.ParamLocationPath, ctx.Param("scanID"), &scanID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter scanID: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetInstancesInstanceIDScanresultsScanIDMalwares(ctx, instanceID, scanID)
	return err
}

// GetInstancesInstanceIDScanresultsScanIDMisconfiguration converts echo context to params.
func (w *ServerInterfaceWrapper) GetInstancesInstanceIDScanresultsScanIDMisconfiguration(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "instanceID" -------------
	var instanceID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "instanceID", runtime.ParamLocationPath, ctx.Param("instanceID"), &instanceID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter instanceID: %s", err))
	}

	// ------------- Path parameter "scanID" -------------
	var scanID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "scanID", runtime.ParamLocationPath, ctx.Param("scanID"), &scanID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter scanID: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetInstancesInstanceIDScanresultsScanIDMisconfiguration(ctx, instanceID, scanID)
	return err
}

// GetInstancesInstanceIDScanresultsScanIDRootkits converts echo context to params.
func (w *ServerInterfaceWrapper) GetInstancesInstanceIDScanresultsScanIDRootkits(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "instanceID" -------------
	var instanceID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "instanceID", runtime.ParamLocationPath, ctx.Param("instanceID"), &instanceID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter instanceID: %s", err))
	}

	// ------------- Path parameter "scanID" -------------
	var scanID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "scanID", runtime.ParamLocationPath, ctx.Param("scanID"), &scanID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter scanID: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetInstancesInstanceIDScanresultsScanIDRootkits(ctx, instanceID, scanID)
	return err
}

// GetInstancesInstanceIDScanresultsScanIDSbom converts echo context to params.
func (w *ServerInterfaceWrapper) GetInstancesInstanceIDScanresultsScanIDSbom(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "instanceID" -------------
	var instanceID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "instanceID", runtime.ParamLocationPath, ctx.Param("instanceID"), &instanceID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter instanceID: %s", err))
	}

	// ------------- Path parameter "scanID" -------------
	var scanID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "scanID", runtime.ParamLocationPath, ctx.Param("scanID"), &scanID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter scanID: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetInstancesInstanceIDScanresultsScanIDSbom(ctx, instanceID, scanID)
	return err
}

// GetInstancesInstanceIDScanresultsScanIDSecrets converts echo context to params.
func (w *ServerInterfaceWrapper) GetInstancesInstanceIDScanresultsScanIDSecrets(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "instanceID" -------------
	var instanceID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "instanceID", runtime.ParamLocationPath, ctx.Param("instanceID"), &instanceID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter instanceID: %s", err))
	}

	// ------------- Path parameter "scanID" -------------
	var scanID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "scanID", runtime.ParamLocationPath, ctx.Param("scanID"), &scanID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter scanID: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetInstancesInstanceIDScanresultsScanIDSecrets(ctx, instanceID, scanID)
	return err
}

// GetInstancesInstanceIDScanresultsScanIDVulnerabilities converts echo context to params.
func (w *ServerInterfaceWrapper) GetInstancesInstanceIDScanresultsScanIDVulnerabilities(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "instanceID" -------------
	var instanceID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "instanceID", runtime.ParamLocationPath, ctx.Param("instanceID"), &instanceID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter instanceID: %s", err))
	}

	// ------------- Path parameter "scanID" -------------
	var scanID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "scanID", runtime.ParamLocationPath, ctx.Param("scanID"), &scanID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter scanID: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetInstancesInstanceIDScanresultsScanIDVulnerabilities(ctx, instanceID, scanID)
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

	router.GET(baseURL+"/instances", wrapper.GetInstances)
	router.POST(baseURL+"/instances", wrapper.PostInstances)
	router.DELETE(baseURL+"/instances/:instanceID", wrapper.DeleteInstancesInstanceID)
	router.GET(baseURL+"/instances/:instanceID", wrapper.GetInstancesInstanceID)
	router.PUT(baseURL+"/instances/:instanceID", wrapper.PutInstancesInstanceID)
	router.GET(baseURL+"/instances/:instanceID/scanresults", wrapper.GetInstancesInstanceIDScanresults)
	router.POST(baseURL+"/instances/:instanceID/scanresults", wrapper.PostInstancesInstanceIDScanresults)
	router.GET(baseURL+"/instances/:instanceID/scanresults/:scanID", wrapper.GetInstancesInstanceIDScanresultsScanID)
	router.GET(baseURL+"/instances/:instanceID/scanresults/:scanID/exploits", wrapper.GetInstancesInstanceIDScanresultsScanIDExploits)
	router.GET(baseURL+"/instances/:instanceID/scanresults/:scanID/malwares", wrapper.GetInstancesInstanceIDScanresultsScanIDMalwares)
	router.GET(baseURL+"/instances/:instanceID/scanresults/:scanID/misconfiguration", wrapper.GetInstancesInstanceIDScanresultsScanIDMisconfiguration)
	router.GET(baseURL+"/instances/:instanceID/scanresults/:scanID/rootkits", wrapper.GetInstancesInstanceIDScanresultsScanIDRootkits)
	router.GET(baseURL+"/instances/:instanceID/scanresults/:scanID/sbom", wrapper.GetInstancesInstanceIDScanresultsScanIDSbom)
	router.GET(baseURL+"/instances/:instanceID/scanresults/:scanID/secrets", wrapper.GetInstancesInstanceIDScanresultsScanIDSecrets)
	router.GET(baseURL+"/instances/:instanceID/scanresults/:scanID/vulnerabilities", wrapper.GetInstancesInstanceIDScanresultsScanIDVulnerabilities)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/+xaX3PauhL/Khrd++gAaXsfLm/cNPcMpyVh7DSdTicPwl5AjS25kpyEk+G7n5Esg21s",
	"xxRIwknfwFrt39/uemU9Yp9HMWfAlMT9RxwTQSJQIMw/yqQizIfhx/Qf7uOYqDl2MCMR4H6ewMECfiZU",
	"QID7SiTgYOnPISJ6p1rEmloqQdkML5cOjskM9EoA0hc0VpRr5mMyA8SSaAIC8SlSc0A/ExAL7KTCsz9W",
	"umHSQi5lCmYgVoI9+leF8BF5oFESIaogkkhxJEAlgjXINnya5EcpS9z/T8/BEWXpn1OnSjHpE1brZ7u4",
	"jY+XmljGnEkwsfQS3wdpfvqcKWBK/yRxHFKfaBd0f0jth8ccz38LmOI+/ld3DZJuuiq7lp9rZaQSix6V",
	"KQmKQEodqqWDv7Bbxu/ZuRBc7E2VQUyb1EhSmQiM0NTXZqPmm9+7AYkBQ3zyA3yF1JwoRKXFBASIMkTC",
	"EPlEgtRYnRIaJgJkBzs4FjwGoWjq+Mz6yjywT1Ip2kFnIU+CseB3NADjIWAaM9/x4KuHb5wyCwefP8Qh",
	"p2rIptwkcEF2wZzHzc00qHx8l4QMBJnQkGacTFZUEtsHRAiyqLbJquj5hG2qGBP/lsxKQpqinTe4lfih",
	"LVKbsmvML5a92uULk5sNBPkwNllUjHmOgQuzusjpiuCCTMK0bk+5iIjSYKdMvX+HK0tMrWeqwXO03qky",
	"dUTCeyK2szRK99RaYtevzPNmI0Y5UtOF1Lyq/al51vemNIS06ugSSSiTyIrDzhb27inn8t5rlXOjomtW",
	"Rezj14F7jh18PXS/eNjBV+7ln4ML7OCvl+4IO9gbf7MU7uDCuxyZP1Vlb0Slz9mUzhJhWsZe69+vhaek",
	"kWwZqNK2fUWsyj+tQjdORbXOEqtaFoAmpcY50ibZW2WplV+bpXb9GoRsXy5cztVtXVfdJ2pEKghXQNwu",
	"1Rpm19uUHzdH2mTvnsCX914rzLlFU7JyMTofXbrfsIM/nbsX55+xgwfj8efh2eBqeKmLxv+Hbn2F8CY8",
	"2pM9WUq0ssUrduaiaEhfXtq+5Bj96+uU7QhPcrvOvc8tMp6b9WrLkpIxkhP+tAtX0dAbwBfwtA88Q5Zt",
	"qngn3dLkJ4LlJVFExKI+Zmc8SSeV8qvVOhBNJGV3N9BmAG0gsT5soCg5rJay0i2G+8v31AwprTppDi+7",
	"p3zOA+2yvjQObzlLrubIbGY2RCmFRIyrOWWzXWfLQk4cor8XBNR3+U2yQwyui5rG+aRG1Qhqmolb+6QN",
	"lpZmxEmdoqgK9dr16CwkgqoFGoyHOh3usrcZfNrpdXqaL4+BkZjiPn7f6XVOcZpyRr9uNjKZfzMwdUAb",
	"l74aBriP/9Bgz4icwjng9+c/pntdZ3NlbTwuFLqFOnslF+qTWW1/WndTOq171+ttdTxWyuQMmEW9P1Op",
	"dMBWaNCVR0dvRu+A6fqrA44IC5B2YQeZDSGwmZqjKJEKTQCF/F6HXSD4mZBQxyJzty5PrVJidSCzkQ0O",
	"VlyRcFP1K/04p7ivuxlKWABi04Lq0491ML5bKTeVuVeKtT02NQtTkoTqBY8tszcUna1rZ5gizWVFUo+5",
	"LGS1dgFI9T8eLPZmReEMaVl0tEb9cgPZpweUXXThBdyv3ITuiUS+AKIg6GiXfej999kUydYRCQWQYIHg",
	"gUrVeWWoOjPeWXkMO/jhxOcBzICdWOicTHiwOLGFTv82LNb9pfu4Ph5cpmkcgoJNaH40z1fgHOY/JJWa",
	"z/6+O21W2Q+bpSaX8R+q1lexZFyhKU9Y8NrimPoWZYoa9Z5s+i8Vgd6zJeHrL+UmUnFSVcmTF4nUa2gX",
	"zweQL3FAcnmTzWLTJOwccTWwVg3Wsjs7V/au9AkT60OuLaqLl9t5QPg6v4eWf9bQohGHLHDQlAtEkIzB",
	"p1NqTlBsFr7MQJM/8t12pinY9XusaRfnliPPs1edAzXNAr6ed8TaEP1CE9YTenh51BzDkJWDeT3K99mm",
	"u4/pJbLlr/drL7uFdsiuvZ/bb4ecMyq+GR1NhXUh1u8BK7T59UD0c9V2S4h18585d8HaecbnrWOu8D34",
	"eMFmcXEY0OW/hu8CulHG562DLn+H7IhBZ3FxINCVPuvvDL4yvzcPwspbJ0eMxpI9h4GlvRu2cy10Mz5v",
	"HYb5C3JHjD6Li8OATk54tCvgPM3jzU8Yq8tyRzxXTHh0IJitrw/uhLTV3bI3DrbcNcsjhpux4jCAq7j6",
	"tQvwrkvs3joAq67qHi0OS1hpCUhzoVfcZQBIRIj7uEtiipc3y78DAAD//8tmWoykOwAA",
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
