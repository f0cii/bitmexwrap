// Code generated by go-swagger; DO NOT EDIT.

package instrument

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/SuperGod/coinex/bitmex/models"
)

// InstrumentGetIndicesReader is a Reader for the InstrumentGetIndices structure.
type InstrumentGetIndicesReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *InstrumentGetIndicesReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewInstrumentGetIndicesOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	case 400:
		result := NewInstrumentGetIndicesBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	case 401:
		result := NewInstrumentGetIndicesUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	case 404:
		result := NewInstrumentGetIndicesNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewInstrumentGetIndicesOK creates a InstrumentGetIndicesOK with default headers values
func NewInstrumentGetIndicesOK() *InstrumentGetIndicesOK {
	return &InstrumentGetIndicesOK{}
}

/*InstrumentGetIndicesOK handles this case with default header values.

Request was successful
*/
type InstrumentGetIndicesOK struct {
	Payload []*models.Instrument
}

func (o *InstrumentGetIndicesOK) Error() string {
	return fmt.Sprintf("[GET /instrument/indices][%d] instrumentGetIndicesOK  %+v", 200, o.Payload)
}

func (o *InstrumentGetIndicesOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewInstrumentGetIndicesBadRequest creates a InstrumentGetIndicesBadRequest with default headers values
func NewInstrumentGetIndicesBadRequest() *InstrumentGetIndicesBadRequest {
	return &InstrumentGetIndicesBadRequest{}
}

/*InstrumentGetIndicesBadRequest handles this case with default header values.

Parameter Error
*/
type InstrumentGetIndicesBadRequest struct {
	Payload *models.Error
}

func (o *InstrumentGetIndicesBadRequest) Error() string {
	return fmt.Sprintf("[GET /instrument/indices][%d] instrumentGetIndicesBadRequest  %+v", 400, o.Payload)
}

func (o *InstrumentGetIndicesBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewInstrumentGetIndicesUnauthorized creates a InstrumentGetIndicesUnauthorized with default headers values
func NewInstrumentGetIndicesUnauthorized() *InstrumentGetIndicesUnauthorized {
	return &InstrumentGetIndicesUnauthorized{}
}

/*InstrumentGetIndicesUnauthorized handles this case with default header values.

Unauthorized
*/
type InstrumentGetIndicesUnauthorized struct {
	Payload *models.Error
}

func (o *InstrumentGetIndicesUnauthorized) Error() string {
	return fmt.Sprintf("[GET /instrument/indices][%d] instrumentGetIndicesUnauthorized  %+v", 401, o.Payload)
}

func (o *InstrumentGetIndicesUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewInstrumentGetIndicesNotFound creates a InstrumentGetIndicesNotFound with default headers values
func NewInstrumentGetIndicesNotFound() *InstrumentGetIndicesNotFound {
	return &InstrumentGetIndicesNotFound{}
}

/*InstrumentGetIndicesNotFound handles this case with default header values.

Not Found
*/
type InstrumentGetIndicesNotFound struct {
	Payload *models.Error
}

func (o *InstrumentGetIndicesNotFound) Error() string {
	return fmt.Sprintf("[GET /instrument/indices][%d] instrumentGetIndicesNotFound  %+v", 404, o.Payload)
}

func (o *InstrumentGetIndicesNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
