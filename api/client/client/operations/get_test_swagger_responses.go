// Code generated by go-swagger; DO NOT EDIT.

package operations

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/strfmt"

	"github.com/openclarity/vmclarity/api/client/models"
)

// GetTestReader is a Reader for the GetTest structure.
type GetTestReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *GetTestReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {
	case 200:
		result := NewGetTestOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil
	default:
		result := NewGetTestDefault(response.Code())
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		if response.Code()/100 == 2 {
			return result, nil
		}
		return nil, result
	}
}

// NewGetTestOK creates a GetTestOK with default headers values
func NewGetTestOK() *GetTestOK {
	return &GetTestOK{}
}

/* GetTestOK describes a response with status code 200, with default header values.

Success
*/
type GetTestOK struct {
	Payload *models.Test
}

func (o *GetTestOK) Error() string {
	return fmt.Sprintf("[GET /test][%d] getTestOK  %+v", 200, o.Payload)
}
func (o *GetTestOK) GetPayload() *models.Test {
	return o.Payload
}

func (o *GetTestOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Test)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewGetTestDefault creates a GetTestDefault with default headers values
func NewGetTestDefault(code int) *GetTestDefault {
	return &GetTestDefault{
		_statusCode: code,
	}
}

/* GetTestDefault describes a response with status code -1, with default header values.

unknown error
*/
type GetTestDefault struct {
	_statusCode int

	Payload *models.APIResponse
}

// Code gets the status code for the get test default response
func (o *GetTestDefault) Code() int {
	return o._statusCode
}

func (o *GetTestDefault) Error() string {
	return fmt.Sprintf("[GET /test][%d] GetTest default  %+v", o._statusCode, o.Payload)
}
func (o *GetTestDefault) GetPayload() *models.APIResponse {
	return o.Payload
}

func (o *GetTestDefault) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.APIResponse)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
