// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/openclarity/vmclarity/api/server/models"
)

// GetPackagesIDOKCode is the HTTP code returned for type GetPackagesIDOK
const GetPackagesIDOKCode int = 200

/*
GetPackagesIDOK Success

swagger:response getPackagesIdOK
*/
type GetPackagesIDOK struct {

	/*
	  In: Body
	*/
	Payload *models.Package `json:"body,omitempty"`
}

// NewGetPackagesIDOK creates GetPackagesIDOK with default headers values
func NewGetPackagesIDOK() *GetPackagesIDOK {

	return &GetPackagesIDOK{}
}

// WithPayload adds the payload to the get packages Id o k response
func (o *GetPackagesIDOK) WithPayload(payload *models.Package) *GetPackagesIDOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get packages Id o k response
func (o *GetPackagesIDOK) SetPayload(payload *models.Package) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetPackagesIDOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// GetPackagesIDNotFoundCode is the HTTP code returned for type GetPackagesIDNotFound
const GetPackagesIDNotFoundCode int = 404

/*
GetPackagesIDNotFound Package ID not found.

swagger:response getPackagesIdNotFound
*/
type GetPackagesIDNotFound struct {
}

// NewGetPackagesIDNotFound creates GetPackagesIDNotFound with default headers values
func NewGetPackagesIDNotFound() *GetPackagesIDNotFound {

	return &GetPackagesIDNotFound{}
}

// WriteResponse to the client
func (o *GetPackagesIDNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(404)
}

/*
GetPackagesIDDefault unknown error

swagger:response getPackagesIdDefault
*/
type GetPackagesIDDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.APIResponse `json:"body,omitempty"`
}

// NewGetPackagesIDDefault creates GetPackagesIDDefault with default headers values
func NewGetPackagesIDDefault(code int) *GetPackagesIDDefault {
	if code <= 0 {
		code = 500
	}

	return &GetPackagesIDDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the get packages ID default response
func (o *GetPackagesIDDefault) WithStatusCode(code int) *GetPackagesIDDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the get packages ID default response
func (o *GetPackagesIDDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the get packages ID default response
func (o *GetPackagesIDDefault) WithPayload(payload *models.APIResponse) *GetPackagesIDDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the get packages ID default response
func (o *GetPackagesIDDefault) SetPayload(payload *models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *GetPackagesIDDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
