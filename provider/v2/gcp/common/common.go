package common

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
  "time"
	
  log "github.com/sirupsen/logrus"
	"google.golang.org/api/googleapi"

	"github.com/openclarity/vmclarity/provider"
)

func HandleGcpRequestError(err error, actionTmpl string, parts ...interface{}) (bool, error) {
	action := fmt.Sprintf(actionTmpl, parts...)

	var gAPIError *googleapi.Error
	if errors.As(err, &gAPIError) {
		sc := gAPIError.Code
		switch {
		case sc >= http.StatusBadRequest && sc < http.StatusInternalServerError:
			// Client errors (BadRequest/Unauthorized/etc.) are Fatal. We
			// also return true to indicate we have NotFound which is a
			// special case in a lot of processing.
			return sc == http.StatusNotFound, provider.FatalErrorf("error from gcp while %s: %w", action, gAPIError)
		default:
			// Everything else is a normal error which can be
			// logged as a failure and then the reconciler will try
			// again on the next loop.
			return false, fmt.Errorf("error from gcp while %s: %w", action, gAPIError)
		}
	} else {
		// Error should be a googleapi.Error
		return false, provider.FatalErrorf("unexpected error from gcp while %s: %w", action, err)
	}
}

func EnsureDeleted(resourceType string, getFunc func() error, deleteFunc func() error, estimateTime time.Duration) error {
	err := getFunc()
	if err != nil {
		notFound, err := HandleGcpRequestError(err, "getting %s", resourceType)
		// NotFound means that the resource has been deleted
		// successfully, all other errors are raised.
		if notFound {
			return nil
		}
		return err
	}

	err = deleteFunc()
	if err != nil {
		_, err := HandleGcpRequestError(err, "deleting %s", resourceType)
		return err
	}

	return provider.RetryableErrorf(estimateTime, "%s delete issued", resourceType)
}

// example: https://www.googleapis.com/compute/v1/projects/gcp-etigcp-nprd-12855/zones/us-central1-c/machineTypes/e2-medium -> returns e2-medium
func GetLastURLPart(str *string) string {
	if str == nil {
		return ""
	}

	urlParsed, err := url.Parse(*str)
	if err != nil {
		log.Error(err)
		return ""
	}

	return path.Base(urlParsed.Path)
}
