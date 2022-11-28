// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// NewGetVulnerabilitiesParams creates a new GetVulnerabilitiesParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewGetVulnerabilitiesParams() *GetVulnerabilitiesParams {
	return &GetVulnerabilitiesParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewGetVulnerabilitiesParamsWithTimeout creates a new GetVulnerabilitiesParams object
// with the ability to set a timeout on a request.
func NewGetVulnerabilitiesParamsWithTimeout(timeout time.Duration) *GetVulnerabilitiesParams {
	return &GetVulnerabilitiesParams{
		timeout: timeout,
	}
}

// NewGetVulnerabilitiesParamsWithContext creates a new GetVulnerabilitiesParams object
// with the ability to set a context for a request.
func NewGetVulnerabilitiesParamsWithContext(ctx context.Context) *GetVulnerabilitiesParams {
	return &GetVulnerabilitiesParams{
		Context: ctx,
	}
}

// NewGetVulnerabilitiesParamsWithHTTPClient creates a new GetVulnerabilitiesParams object
// with the ability to set a custom HTTPClient for a request.
func NewGetVulnerabilitiesParamsWithHTTPClient(client *http.Client) *GetVulnerabilitiesParams {
	return &GetVulnerabilitiesParams{
		HTTPClient: client,
	}
}

/*
GetVulnerabilitiesParams contains all the parameters to send to the API endpoint

	for the get vulnerabilities operation.

	Typically these are written to a http.Request.
*/
type GetVulnerabilitiesParams struct {

	/* Page.

	   Page number of the query
	*/
	Page int64

	/* PageSize.

	   Maximum items to return
	*/
	PageSize int64

	/* SortKey.

	   Sort key
	*/
	SortKey string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the get vulnerabilities params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetVulnerabilitiesParams) WithDefaults() *GetVulnerabilitiesParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the get vulnerabilities params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetVulnerabilitiesParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the get vulnerabilities params
func (o *GetVulnerabilitiesParams) WithTimeout(timeout time.Duration) *GetVulnerabilitiesParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get vulnerabilities params
func (o *GetVulnerabilitiesParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get vulnerabilities params
func (o *GetVulnerabilitiesParams) WithContext(ctx context.Context) *GetVulnerabilitiesParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get vulnerabilities params
func (o *GetVulnerabilitiesParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get vulnerabilities params
func (o *GetVulnerabilitiesParams) WithHTTPClient(client *http.Client) *GetVulnerabilitiesParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get vulnerabilities params
func (o *GetVulnerabilitiesParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithPage adds the page to the get vulnerabilities params
func (o *GetVulnerabilitiesParams) WithPage(page int64) *GetVulnerabilitiesParams {
	o.SetPage(page)
	return o
}

// SetPage adds the page to the get vulnerabilities params
func (o *GetVulnerabilitiesParams) SetPage(page int64) {
	o.Page = page
}

// WithPageSize adds the pageSize to the get vulnerabilities params
func (o *GetVulnerabilitiesParams) WithPageSize(pageSize int64) *GetVulnerabilitiesParams {
	o.SetPageSize(pageSize)
	return o
}

// SetPageSize adds the pageSize to the get vulnerabilities params
func (o *GetVulnerabilitiesParams) SetPageSize(pageSize int64) {
	o.PageSize = pageSize
}

// WithSortKey adds the sortKey to the get vulnerabilities params
func (o *GetVulnerabilitiesParams) WithSortKey(sortKey string) *GetVulnerabilitiesParams {
	o.SetSortKey(sortKey)
	return o
}

// SetSortKey adds the sortKey to the get vulnerabilities params
func (o *GetVulnerabilitiesParams) SetSortKey(sortKey string) {
	o.SortKey = sortKey
}

// WriteToRequest writes these params to a swagger request
func (o *GetVulnerabilitiesParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// query param page
	qrPage := o.Page
	qPage := swag.FormatInt64(qrPage)
	if qPage != "" {

		if err := r.SetQueryParam("page", qPage); err != nil {
			return err
		}
	}

	// query param pageSize
	qrPageSize := o.PageSize
	qPageSize := swag.FormatInt64(qrPageSize)
	if qPageSize != "" {

		if err := r.SetQueryParam("pageSize", qPageSize); err != nil {
			return err
		}
	}

	// query param sortKey
	qrSortKey := o.SortKey
	qSortKey := qrSortKey
	if qSortKey != "" {

		if err := r.SetQueryParam("sortKey", qSortKey); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
