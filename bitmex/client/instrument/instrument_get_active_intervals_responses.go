// Code generated by go-swagger; DO NOT EDIT.

package instrument

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/frankrap/bitmexwrap/bitmex/models"
)

// InstrumentGetActiveIntervalsReader is a Reader for the InstrumentGetActiveIntervals structure.
type InstrumentGetActiveIntervalsReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *InstrumentGetActiveIntervalsReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewInstrumentGetActiveIntervalsOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	case 400:
		result := NewInstrumentGetActiveIntervalsBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	case 401:
		result := NewInstrumentGetActiveIntervalsUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	case 404:
		result := NewInstrumentGetActiveIntervalsNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewInstrumentGetActiveIntervalsOK creates a InstrumentGetActiveIntervalsOK with default headers values
func NewInstrumentGetActiveIntervalsOK() *InstrumentGetActiveIntervalsOK {
	return &InstrumentGetActiveIntervalsOK{}
}

/*InstrumentGetActiveIntervalsOK handles this case with default header values.

Request was successful
*/
type InstrumentGetActiveIntervalsOK struct {
	Payload *models.InstrumentInterval
}

func (o *InstrumentGetActiveIntervalsOK) Error() string {
	return fmt.Sprintf("[GET /instrument/activeIntervals][%d] instrumentGetActiveIntervalsOK  %+v", 200, *o.Payload)
}

func (o *InstrumentGetActiveIntervalsOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.InstrumentInterval)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewInstrumentGetActiveIntervalsBadRequest creates a InstrumentGetActiveIntervalsBadRequest with default headers values
func NewInstrumentGetActiveIntervalsBadRequest() *InstrumentGetActiveIntervalsBadRequest {
	return &InstrumentGetActiveIntervalsBadRequest{}
}

/*InstrumentGetActiveIntervalsBadRequest handles this case with default header values.

Parameter Error
*/
type InstrumentGetActiveIntervalsBadRequest struct {
	Payload *models.Error
}

func (o *InstrumentGetActiveIntervalsBadRequest) Error() string {
	return fmt.Sprintf("[GET /instrument/activeIntervals][%d] instrumentGetActiveIntervalsBadRequest  %+v", 400, *o.Payload)
}

func (o *InstrumentGetActiveIntervalsBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewInstrumentGetActiveIntervalsUnauthorized creates a InstrumentGetActiveIntervalsUnauthorized with default headers values
func NewInstrumentGetActiveIntervalsUnauthorized() *InstrumentGetActiveIntervalsUnauthorized {
	return &InstrumentGetActiveIntervalsUnauthorized{}
}

/*InstrumentGetActiveIntervalsUnauthorized handles this case with default header values.

Unauthorized
*/
type InstrumentGetActiveIntervalsUnauthorized struct {
	Payload *models.Error
}

func (o *InstrumentGetActiveIntervalsUnauthorized) Error() string {
	return fmt.Sprintf("[GET /instrument/activeIntervals][%d] instrumentGetActiveIntervalsUnauthorized  %+v", 401, o.Payload)
}

func (o *InstrumentGetActiveIntervalsUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewInstrumentGetActiveIntervalsNotFound creates a InstrumentGetActiveIntervalsNotFound with default headers values
func NewInstrumentGetActiveIntervalsNotFound() *InstrumentGetActiveIntervalsNotFound {
	return &InstrumentGetActiveIntervalsNotFound{}
}

/*InstrumentGetActiveIntervalsNotFound handles this case with default header values.

Not Found
*/
type InstrumentGetActiveIntervalsNotFound struct {
	Payload *models.Error
}

func (o *InstrumentGetActiveIntervalsNotFound) Error() string {
	return fmt.Sprintf("[GET /instrument/activeIntervals][%d] instrumentGetActiveIntervalsNotFound  %+v", 404, o.Payload)
}

func (o *InstrumentGetActiveIntervalsNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
