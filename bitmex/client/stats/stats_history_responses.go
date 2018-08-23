// Code generated by go-swagger; DO NOT EDIT.

package stats

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/SuperGod/coinex/bitmex/models"
)

// StatsHistoryReader is a Reader for the StatsHistory structure.
type StatsHistoryReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *StatsHistoryReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewStatsHistoryOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	case 400:
		result := NewStatsHistoryBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	case 401:
		result := NewStatsHistoryUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	case 404:
		result := NewStatsHistoryNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewStatsHistoryOK creates a StatsHistoryOK with default headers values
func NewStatsHistoryOK() *StatsHistoryOK {
	return &StatsHistoryOK{}
}

/*StatsHistoryOK handles this case with default header values.

Request was successful
*/
type StatsHistoryOK struct {
	Payload []*models.StatsHistory
}

func (o *StatsHistoryOK) Error() string {
	return fmt.Sprintf("[GET /stats/history][%d] statsHistoryOK  %+v", 200, o.Payload)
}

func (o *StatsHistoryOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewStatsHistoryBadRequest creates a StatsHistoryBadRequest with default headers values
func NewStatsHistoryBadRequest() *StatsHistoryBadRequest {
	return &StatsHistoryBadRequest{}
}

/*StatsHistoryBadRequest handles this case with default header values.

Parameter Error
*/
type StatsHistoryBadRequest struct {
	Payload *models.Error
}

func (o *StatsHistoryBadRequest) Error() string {
	return fmt.Sprintf("[GET /stats/history][%d] statsHistoryBadRequest  %+v", 400, o.Payload)
}

func (o *StatsHistoryBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewStatsHistoryUnauthorized creates a StatsHistoryUnauthorized with default headers values
func NewStatsHistoryUnauthorized() *StatsHistoryUnauthorized {
	return &StatsHistoryUnauthorized{}
}

/*StatsHistoryUnauthorized handles this case with default header values.

Unauthorized
*/
type StatsHistoryUnauthorized struct {
	Payload *models.Error
}

func (o *StatsHistoryUnauthorized) Error() string {
	return fmt.Sprintf("[GET /stats/history][%d] statsHistoryUnauthorized  %+v", 401, o.Payload)
}

func (o *StatsHistoryUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewStatsHistoryNotFound creates a StatsHistoryNotFound with default headers values
func NewStatsHistoryNotFound() *StatsHistoryNotFound {
	return &StatsHistoryNotFound{}
}

/*StatsHistoryNotFound handles this case with default header values.

Not Found
*/
type StatsHistoryNotFound struct {
	Payload *models.Error
}

func (o *StatsHistoryNotFound) Error() string {
	return fmt.Sprintf("[GET /stats/history][%d] statsHistoryNotFound  %+v", 404, o.Payload)
}

func (o *StatsHistoryNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
