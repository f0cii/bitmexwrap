// Code generated by go-swagger; DO NOT EDIT.

package trade

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/SuperGod/coinex/bitmex/models"
)

// TradeGetReader is a Reader for the TradeGet structure.
type TradeGetReader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *TradeGetReader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewTradeGetOK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	case 400:
		result := NewTradeGetBadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	case 401:
		result := NewTradeGetUnauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	case 404:
		result := NewTradeGetNotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewTradeGetOK creates a TradeGetOK with default headers values
func NewTradeGetOK() *TradeGetOK {
	return &TradeGetOK{}
}

/*TradeGetOK handles this case with default header values.

Request was successful
*/
type TradeGetOK struct {
	Payload []*models.Trade
}

func (o *TradeGetOK) Error() string {
	return fmt.Sprintf("[GET /trade][%d] tradeGetOK  %+v", 200, o.Payload)
}

func (o *TradeGetOK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewTradeGetBadRequest creates a TradeGetBadRequest with default headers values
func NewTradeGetBadRequest() *TradeGetBadRequest {
	return &TradeGetBadRequest{}
}

/*TradeGetBadRequest handles this case with default header values.

Parameter Error
*/
type TradeGetBadRequest struct {
	Payload *models.Error
}

func (o *TradeGetBadRequest) Error() string {
	return fmt.Sprintf("[GET /trade][%d] tradeGetBadRequest  %+v", 400, o.Payload)
}

func (o *TradeGetBadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewTradeGetUnauthorized creates a TradeGetUnauthorized with default headers values
func NewTradeGetUnauthorized() *TradeGetUnauthorized {
	return &TradeGetUnauthorized{}
}

/*TradeGetUnauthorized handles this case with default header values.

Unauthorized
*/
type TradeGetUnauthorized struct {
	Payload *models.Error
}

func (o *TradeGetUnauthorized) Error() string {
	return fmt.Sprintf("[GET /trade][%d] tradeGetUnauthorized  %+v", 401, o.Payload)
}

func (o *TradeGetUnauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewTradeGetNotFound creates a TradeGetNotFound with default headers values
func NewTradeGetNotFound() *TradeGetNotFound {
	return &TradeGetNotFound{}
}

/*TradeGetNotFound handles this case with default header values.

Not Found
*/
type TradeGetNotFound struct {
	Payload *models.Error
}

func (o *TradeGetNotFound) Error() string {
	return fmt.Sprintf("[GET /trade][%d] tradeGetNotFound  %+v", 404, o.Payload)
}

func (o *TradeGetNotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
