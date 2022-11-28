// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// NewGetInstancesParams creates a new GetInstancesParams object
//
// There are no default values defined in the spec.
func NewGetInstancesParams() GetInstancesParams {

	return GetInstancesParams{}
}

// GetInstancesParams contains all the bound params for the get instances operation
// typically these are obtained from a http.Request
//
// swagger:parameters GetInstances
type GetInstancesParams struct {

	// HTTP Request Object
	HTTPRequest *http.Request `json:"-"`

	/*Page number of the query
	  Required: true
	  In: query
	*/
	Page int64
	/*Maximum items to return
	  Required: true
	  Maximum: 50
	  Minimum: 1
	  In: query
	*/
	PageSize int64
	/*Sort key
	  Required: true
	  In: query
	*/
	SortKey string
}

// BindRequest both binds and validates a request, it assumes that complex things implement a Validatable(strfmt.Registry) error interface
// for simple values it will use straight method calls.
//
// To ensure default values, the struct must have been initialized with NewGetInstancesParams() beforehand.
func (o *GetInstancesParams) BindRequest(r *http.Request, route *middleware.MatchedRoute) error {
	var res []error

	o.HTTPRequest = r

	qs := runtime.Values(r.URL.Query())

	qPage, qhkPage, _ := qs.GetOK("page")
	if err := o.bindPage(qPage, qhkPage, route.Formats); err != nil {
		res = append(res, err)
	}

	qPageSize, qhkPageSize, _ := qs.GetOK("pageSize")
	if err := o.bindPageSize(qPageSize, qhkPageSize, route.Formats); err != nil {
		res = append(res, err)
	}

	qSortKey, qhkSortKey, _ := qs.GetOK("sortKey")
	if err := o.bindSortKey(qSortKey, qhkSortKey, route.Formats); err != nil {
		res = append(res, err)
	}
	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// bindPage binds and validates parameter Page from query.
func (o *GetInstancesParams) bindPage(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("page", "query", rawData)
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// AllowEmptyValue: false

	if err := validate.RequiredString("page", "query", raw); err != nil {
		return err
	}

	value, err := swag.ConvertInt64(raw)
	if err != nil {
		return errors.InvalidType("page", "query", "int64", raw)
	}
	o.Page = value

	return nil
}

// bindPageSize binds and validates parameter PageSize from query.
func (o *GetInstancesParams) bindPageSize(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("pageSize", "query", rawData)
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// AllowEmptyValue: false

	if err := validate.RequiredString("pageSize", "query", raw); err != nil {
		return err
	}

	value, err := swag.ConvertInt64(raw)
	if err != nil {
		return errors.InvalidType("pageSize", "query", "int64", raw)
	}
	o.PageSize = value

	if err := o.validatePageSize(formats); err != nil {
		return err
	}

	return nil
}

// validatePageSize carries on validations for parameter PageSize
func (o *GetInstancesParams) validatePageSize(formats strfmt.Registry) error {

	if err := validate.MinimumInt("pageSize", "query", o.PageSize, 1, false); err != nil {
		return err
	}

	if err := validate.MaximumInt("pageSize", "query", o.PageSize, 50, false); err != nil {
		return err
	}

	return nil
}

// bindSortKey binds and validates parameter SortKey from query.
func (o *GetInstancesParams) bindSortKey(rawData []string, hasKey bool, formats strfmt.Registry) error {
	if !hasKey {
		return errors.Required("sortKey", "query", rawData)
	}
	var raw string
	if len(rawData) > 0 {
		raw = rawData[len(rawData)-1]
	}

	// Required: true
	// AllowEmptyValue: false

	if err := validate.RequiredString("sortKey", "query", raw); err != nil {
		return err
	}
	o.SortKey = raw

	if err := o.validateSortKey(formats); err != nil {
		return err
	}

	return nil
}

// validateSortKey carries on validations for parameter SortKey
func (o *GetInstancesParams) validateSortKey(formats strfmt.Registry) error {

	if err := validate.EnumCase("sortKey", "query", o.SortKey, []interface{}{"instanceName", "instanceProvider", "instanceRegion", "instanceID"}, true); err != nil {
		return err
	}

	return nil
}
