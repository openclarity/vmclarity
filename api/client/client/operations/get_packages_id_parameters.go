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
)

// NewGetPackagesIDParams creates a new GetPackagesIDParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewGetPackagesIDParams() *GetPackagesIDParams {
	return &GetPackagesIDParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewGetPackagesIDParamsWithTimeout creates a new GetPackagesIDParams object
// with the ability to set a timeout on a request.
func NewGetPackagesIDParamsWithTimeout(timeout time.Duration) *GetPackagesIDParams {
	return &GetPackagesIDParams{
		timeout: timeout,
	}
}

// NewGetPackagesIDParamsWithContext creates a new GetPackagesIDParams object
// with the ability to set a context for a request.
func NewGetPackagesIDParamsWithContext(ctx context.Context) *GetPackagesIDParams {
	return &GetPackagesIDParams{
		Context: ctx,
	}
}

// NewGetPackagesIDParamsWithHTTPClient creates a new GetPackagesIDParams object
// with the ability to set a custom HTTPClient for a request.
func NewGetPackagesIDParamsWithHTTPClient(client *http.Client) *GetPackagesIDParams {
	return &GetPackagesIDParams{
		HTTPClient: client,
	}
}

/*
GetPackagesIDParams contains all the parameters to send to the API endpoint

	for the get packages ID operation.

	Typically these are written to a http.Request.
*/
type GetPackagesIDParams struct {

	// ID.
	ID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the get packages ID params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetPackagesIDParams) WithDefaults() *GetPackagesIDParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the get packages ID params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetPackagesIDParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the get packages ID params
func (o *GetPackagesIDParams) WithTimeout(timeout time.Duration) *GetPackagesIDParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get packages ID params
func (o *GetPackagesIDParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get packages ID params
func (o *GetPackagesIDParams) WithContext(ctx context.Context) *GetPackagesIDParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get packages ID params
func (o *GetPackagesIDParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get packages ID params
func (o *GetPackagesIDParams) WithHTTPClient(client *http.Client) *GetPackagesIDParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get packages ID params
func (o *GetPackagesIDParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithID adds the id to the get packages ID params
func (o *GetPackagesIDParams) WithID(id string) *GetPackagesIDParams {
	o.SetID(id)
	return o
}

// SetID adds the id to the get packages ID params
func (o *GetPackagesIDParams) SetID(id string) {
	o.ID = id
}

// WriteToRequest writes these params to a swagger request
func (o *GetPackagesIDParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param id
	if err := r.SetPathParam("id", o.ID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
