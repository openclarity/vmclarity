// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/openclarity/vmclarity/api/server/restapi/operations"
)

//go:generate swagger generate server --target ../../server --name VMClarityAPIs --spec ../../swagger.yaml --principal interface{}

func configureFlags(api *operations.VMClarityAPIsAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.VMClarityAPIsAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	if api.DeleteInstancesIDHandler == nil {
		api.DeleteInstancesIDHandler = operations.DeleteInstancesIDHandlerFunc(func(params operations.DeleteInstancesIDParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.DeleteInstancesID has not yet been implemented")
		})
	}
	if api.GetInstancesHandler == nil {
		api.GetInstancesHandler = operations.GetInstancesHandlerFunc(func(params operations.GetInstancesParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetInstances has not yet been implemented")
		})
	}
	if api.GetInstancesIDHandler == nil {
		api.GetInstancesIDHandler = operations.GetInstancesIDHandlerFunc(func(params operations.GetInstancesIDParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetInstancesID has not yet been implemented")
		})
	}
	if api.GetPackagesHandler == nil {
		api.GetPackagesHandler = operations.GetPackagesHandlerFunc(func(params operations.GetPackagesParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetPackages has not yet been implemented")
		})
	}
	if api.GetPackagesIDHandler == nil {
		api.GetPackagesIDHandler = operations.GetPackagesIDHandlerFunc(func(params operations.GetPackagesIDParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetPackagesID has not yet been implemented")
		})
	}
	if api.GetVulnerabilitiesHandler == nil {
		api.GetVulnerabilitiesHandler = operations.GetVulnerabilitiesHandlerFunc(func(params operations.GetVulnerabilitiesParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetVulnerabilities has not yet been implemented")
		})
	}
	if api.GetVulnerabilitiesVulIDPkgIDHandler == nil {
		api.GetVulnerabilitiesVulIDPkgIDHandler = operations.GetVulnerabilitiesVulIDPkgIDHandlerFunc(func(params operations.GetVulnerabilitiesVulIDPkgIDParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.GetVulnerabilitiesVulIDPkgID has not yet been implemented")
		})
	}
	if api.PostInstancesHandler == nil {
		api.PostInstancesHandler = operations.PostInstancesHandlerFunc(func(params operations.PostInstancesParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.PostInstances has not yet been implemented")
		})
	}
	if api.PostInstancesContentAnalysisIDHandler == nil {
		api.PostInstancesContentAnalysisIDHandler = operations.PostInstancesContentAnalysisIDHandlerFunc(func(params operations.PostInstancesContentAnalysisIDParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.PostInstancesContentAnalysisID has not yet been implemented")
		})
	}
	if api.PostInstancesExploitScanIDHandler == nil {
		api.PostInstancesExploitScanIDHandler = operations.PostInstancesExploitScanIDHandlerFunc(func(params operations.PostInstancesExploitScanIDParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.PostInstancesExploitScanID has not yet been implemented")
		})
	}
	if api.PostInstancesMalewareScanIDHandler == nil {
		api.PostInstancesMalewareScanIDHandler = operations.PostInstancesMalewareScanIDHandlerFunc(func(params operations.PostInstancesMalewareScanIDParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.PostInstancesMalewareScanID has not yet been implemented")
		})
	}
	if api.PostInstancesMisconfigurationScanIDHandler == nil {
		api.PostInstancesMisconfigurationScanIDHandler = operations.PostInstancesMisconfigurationScanIDHandlerFunc(func(params operations.PostInstancesMisconfigurationScanIDParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.PostInstancesMisconfigurationScanID has not yet been implemented")
		})
	}
	if api.PostInstancesRootkitScanIDHandler == nil {
		api.PostInstancesRootkitScanIDHandler = operations.PostInstancesRootkitScanIDHandlerFunc(func(params operations.PostInstancesRootkitScanIDParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.PostInstancesRootkitScanID has not yet been implemented")
		})
	}
	if api.PostInstancesSecretScanIDHandler == nil {
		api.PostInstancesSecretScanIDHandler = operations.PostInstancesSecretScanIDHandlerFunc(func(params operations.PostInstancesSecretScanIDParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.PostInstancesSecretScanID has not yet been implemented")
		})
	}
	if api.PostInstancesVulnerabilityScanIDHandler == nil {
		api.PostInstancesVulnerabilityScanIDHandler = operations.PostInstancesVulnerabilityScanIDHandlerFunc(func(params operations.PostInstancesVulnerabilityScanIDParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.PostInstancesVulnerabilityScanID has not yet been implemented")
		})
	}
	if api.PutInstancesIDHandler == nil {
		api.PutInstancesIDHandler = operations.PutInstancesIDHandlerFunc(func(params operations.PutInstancesIDParams) middleware.Responder {
			return middleware.NotImplemented("operation operations.PutInstancesID has not yet been implemented")
		})
	}

	api.PreServerShutdown = func() {}

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix".
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation.
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics.
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
