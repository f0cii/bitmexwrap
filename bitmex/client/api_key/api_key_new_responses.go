// Code generated by go-swagger; DO NOT EDIT.

package api_key

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/SuperGod/coinex/bitmex/models"
)

// APIKeyNewReader is a Reader for the APIKeyNew structure.
type APIKeyNewReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *APIKeyNewReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewAPIKeyNewOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	case 400:
		result := NewAPIKeyNewBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	case 401:
		result := NewAPIKeyNewUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	case 404:
		result := NewAPIKeyNewNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewAPIKeyNewOK creates a APIKeyNewOK with default headers values
func NewAPIKeyNewOK() *APIKeyNewOK {
	return &APIKeyNewOK{}
}

/*APIKeyNewOK handles this case with default header values.

Request was successful
*/
type APIKeyNewOK struct {
	Payload *models.APIKey
}

func (o *APIKeyNewOK) Error() string {
	return fmt.Sprintf("[POST /apiKey][%d] apiKeyNewOK  %+v", 200, o.Payload)
}

func (o *APIKeyNewOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.APIKey)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewAPIKeyNewBadRequest creates a APIKeyNewBadRequest with default headers values
func NewAPIKeyNewBadRequest() *APIKeyNewBadRequest {
	return &APIKeyNewBadRequest{}
}

/*APIKeyNewBadRequest handles this case with default header values.

Parameter Error
*/
type APIKeyNewBadRequest struct {
	Payload *models.Error
}

func (o *APIKeyNewBadRequest) Error() string {
	return fmt.Sprintf("[POST /apiKey][%d] apiKeyNewBadRequest  %+v", 400, o.Payload)
}

func (o *APIKeyNewBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewAPIKeyNewUnauthorized creates a APIKeyNewUnauthorized with default headers values
func NewAPIKeyNewUnauthorized() *APIKeyNewUnauthorized {
	return &APIKeyNewUnauthorized{}
}

/*APIKeyNewUnauthorized handles this case with default header values.

Unauthorized
*/
type APIKeyNewUnauthorized struct {
	Payload *models.Error
}

func (o *APIKeyNewUnauthorized) Error() string {
	return fmt.Sprintf("[POST /apiKey][%d] apiKeyNewUnauthorized  %+v", 401, o.Payload)
}

func (o *APIKeyNewUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewAPIKeyNewNotFound creates a APIKeyNewNotFound with default headers values
func NewAPIKeyNewNotFound() *APIKeyNewNotFound {
	return &APIKeyNewNotFound{}
}

/*APIKeyNewNotFound handles this case with default header values.

Not Found
*/
type APIKeyNewNotFound struct {
	Payload *models.Error
}

func (o *APIKeyNewNotFound) Error() string {
	return fmt.Sprintf("[POST /apiKey][%d] apiKeyNewNotFound  %+v", 404, o.Payload)
}

func (o *APIKeyNewNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
