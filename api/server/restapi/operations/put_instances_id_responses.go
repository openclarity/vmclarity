// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"net/http"

	"github.com/go-openapi/runtime"

	"github.com/openclarity/vmclarity/api/server/models"
)

// PutInstancesIDOKCode is the HTTP code returned for type PutInstancesIDOK
const PutInstancesIDOKCode int = 200

/*
PutInstancesIDOK Update Instance successful.

swagger:response putInstancesIdOK
*/
type PutInstancesIDOK struct {

	/*
	  In: Body
	*/
	Payload *models.InstanceInfo `json:"body,omitempty"`
}

// NewPutInstancesIDOK creates PutInstancesIDOK with default headers values
func NewPutInstancesIDOK() *PutInstancesIDOK {

	return &PutInstancesIDOK{}
}

// WithPayload adds the payload to the put instances Id o k response
func (o *PutInstancesIDOK) WithPayload(payload *models.InstanceInfo) *PutInstancesIDOK {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the put instances Id o k response
func (o *PutInstancesIDOK) SetPayload(payload *models.InstanceInfo) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PutInstancesIDOK) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(200)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}

// PutInstancesIDNotFoundCode is the HTTP code returned for type PutInstancesIDNotFound
const PutInstancesIDNotFoundCode int = 404

/*
PutInstancesIDNotFound Instance not found.

swagger:response putInstancesIdNotFound
*/
type PutInstancesIDNotFound struct {
}

// NewPutInstancesIDNotFound creates PutInstancesIDNotFound with default headers values
func NewPutInstancesIDNotFound() *PutInstancesIDNotFound {

	return &PutInstancesIDNotFound{}
}

// WriteResponse to the client
func (o *PutInstancesIDNotFound) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.Header().Del(runtime.HeaderContentType) //Remove Content-Type on empty responses

	rw.WriteHeader(404)
}

/*
PutInstancesIDDefault unknown error

swagger:response putInstancesIdDefault
*/
type PutInstancesIDDefault struct {
	_statusCode int

	/*
	  In: Body
	*/
	Payload *models.APIResponse `json:"body,omitempty"`
}

// NewPutInstancesIDDefault creates PutInstancesIDDefault with default headers values
func NewPutInstancesIDDefault(code int) *PutInstancesIDDefault {
	if code <= 0 {
		code = 500
	}

	return &PutInstancesIDDefault{
		_statusCode: code,
	}
}

// WithStatusCode adds the status to the put instances ID default response
func (o *PutInstancesIDDefault) WithStatusCode(code int) *PutInstancesIDDefault {
	o._statusCode = code
	return o
}

// SetStatusCode sets the status to the put instances ID default response
func (o *PutInstancesIDDefault) SetStatusCode(code int) {
	o._statusCode = code
}

// WithPayload adds the payload to the put instances ID default response
func (o *PutInstancesIDDefault) WithPayload(payload *models.APIResponse) *PutInstancesIDDefault {
	o.Payload = payload
	return o
}

// SetPayload sets the payload to the put instances ID default response
func (o *PutInstancesIDDefault) SetPayload(payload *models.APIResponse) {
	o.Payload = payload
}

// WriteResponse to the client
func (o *PutInstancesIDDefault) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {

	rw.WriteHeader(o._statusCode)
	if o.Payload != nil {
		payload := o.Payload
		if err := producer.Produce(rw, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	}
}
