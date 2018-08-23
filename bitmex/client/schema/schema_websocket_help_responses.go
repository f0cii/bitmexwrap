// Code generated by go-swagger; DO NOT EDIT.

package schema

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/SuperGod/coinex/bitmex/models"
)

// SchemaWebsocketHelpReader is a Reader for the SchemaWebsocketHelp structure.
type SchemaWebsocketHelpReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *SchemaWebsocketHelpReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewSchemaWebsocketHelpOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	case 400:
		result := NewSchemaWebsocketHelpBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	case 401:
		result := NewSchemaWebsocketHelpUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	case 404:
		result := NewSchemaWebsocketHelpNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewSchemaWebsocketHelpOK creates a SchemaWebsocketHelpOK with default headers values
func NewSchemaWebsocketHelpOK() *SchemaWebsocketHelpOK {
	return &SchemaWebsocketHelpOK{}
}

/*SchemaWebsocketHelpOK handles this case with default header values.

Request was successful
*/
type SchemaWebsocketHelpOK struct {
	Payload interface{}
}

func (o *SchemaWebsocketHelpOK) Error() string {
	return fmt.Sprintf("[GET /schema/websocketHelp][%d] schemaWebsocketHelpOK  %+v", 200, o.Payload)
}

func (o *SchemaWebsocketHelpOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewSchemaWebsocketHelpBadRequest creates a SchemaWebsocketHelpBadRequest with default headers values
func NewSchemaWebsocketHelpBadRequest() *SchemaWebsocketHelpBadRequest {
	return &SchemaWebsocketHelpBadRequest{}
}

/*SchemaWebsocketHelpBadRequest handles this case with default header values.

Parameter Error
*/
type SchemaWebsocketHelpBadRequest struct {
	Payload *models.Error
}

func (o *SchemaWebsocketHelpBadRequest) Error() string {
	return fmt.Sprintf("[GET /schema/websocketHelp][%d] schemaWebsocketHelpBadRequest  %+v", 400, o.Payload)
}

func (o *SchemaWebsocketHelpBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewSchemaWebsocketHelpUnauthorized creates a SchemaWebsocketHelpUnauthorized with default headers values
func NewSchemaWebsocketHelpUnauthorized() *SchemaWebsocketHelpUnauthorized {
	return &SchemaWebsocketHelpUnauthorized{}
}

/*SchemaWebsocketHelpUnauthorized handles this case with default header values.

Unauthorized
*/
type SchemaWebsocketHelpUnauthorized struct {
	Payload *models.Error
}

func (o *SchemaWebsocketHelpUnauthorized) Error() string {
	return fmt.Sprintf("[GET /schema/websocketHelp][%d] schemaWebsocketHelpUnauthorized  %+v", 401, o.Payload)
}

func (o *SchemaWebsocketHelpUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewSchemaWebsocketHelpNotFound creates a SchemaWebsocketHelpNotFound with default headers values
func NewSchemaWebsocketHelpNotFound() *SchemaWebsocketHelpNotFound {
	return &SchemaWebsocketHelpNotFound{}
}

/*SchemaWebsocketHelpNotFound handles this case with default header values.

Not Found
*/
type SchemaWebsocketHelpNotFound struct {
	Payload *models.Error
}

func (o *SchemaWebsocketHelpNotFound) Error() string {
	return fmt.Sprintf("[GET /schema/websocketHelp][%d] schemaWebsocketHelpNotFound  %+v", 404, o.Payload)
}

func (o *SchemaWebsocketHelpNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
