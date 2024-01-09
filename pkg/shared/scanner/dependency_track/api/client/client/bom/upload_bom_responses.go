// Code generated by go-swagger; DO NOT EDIT.

package bom

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"
)

// UploadBomReader is a Reader for the UploadBom structure.
type UploadBomReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *UploadBomReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewUploadBomOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	case 401:
		result := NewUploadBomUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 403:
		result := NewUploadBomForbidden()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	case 404:
		result := NewUploadBomNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result
	default:
		return nil, runtime.NewAPIError("response status code does not match any response statuses defined for this endpoint in the swagger spec", response, response.Code())
	}
}

// NewUploadBomOK creates a UploadBomOK with default headers values
func NewUploadBomOK() *UploadBomOK {
	return &UploadBomOK{}
}

/*
UploadBomOK describes a response with status code 200, with default header values.

Successful
*/
type UploadBomOK struct {
}

func (o *UploadBomOK) Error() string {
	return fmt.Sprintf("[POST /v1/bom][%d] uploadBomOK ", 200)
}

func (o *UploadBomOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewUploadBomUnauthorized creates a UploadBomUnauthorized with default headers values
func NewUploadBomUnauthorized() *UploadBomUnauthorized {
	return &UploadBomUnauthorized{}
}

/*
UploadBomUnauthorized describes a response with status code 401, with default header values.

Unauthorized
*/
type UploadBomUnauthorized struct {
}

func (o *UploadBomUnauthorized) Error() string {
	return fmt.Sprintf("[POST /v1/bom][%d] uploadBomUnauthorized ", 401)
}

func (o *UploadBomUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewUploadBomForbidden creates a UploadBomForbidden with default headers values
func NewUploadBomForbidden() *UploadBomForbidden {
	return &UploadBomForbidden{}
}

/*
UploadBomForbidden describes a response with status code 403, with default header values.

Access to the specified project is forbidden
*/
type UploadBomForbidden struct {
}

func (o *UploadBomForbidden) Error() string {
	return fmt.Sprintf("[POST /v1/bom][%d] uploadBomForbidden ", 403)
}

func (o *UploadBomForbidden) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}

// NewUploadBomNotFound creates a UploadBomNotFound with default headers values
func NewUploadBomNotFound() *UploadBomNotFound {
	return &UploadBomNotFound{}
}

/*
UploadBomNotFound describes a response with status code 404, with default header values.

The project could not be found
*/
type UploadBomNotFound struct {
}

func (o *UploadBomNotFound) Error() string {
	return fmt.Sprintf("[POST /v1/bom][%d] uploadBomNotFound ", 404)
}

func (o *UploadBomNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	return nil
}
