// Code generated by go-swagger; DO NOT EDIT.

package order_book

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"fmt"
	"io"

	"github.com/go-openapi/runtime"

	strfmt "github.com/go-openapi/strfmt"

	models "github.com/SuperGod/coinex/bitmex/models"
)

// OrderBookGetL2Reader is a Reader for the OrderBookGetL2 structure.
type OrderBookGetL2Reader struct {
	formats strfmt.Registry
}

// ReadResponse reads a server response into the received o.
func (o *OrderBookGetL2Reader) ReadResponse(response runtime.ClientResponse, consumer runtime.Consumer) (interface{}, error) {
	switch response.Code() {

	case 200:
		result := NewOrderBookGetL2OK()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return result, nil

	case 400:
		result := NewOrderBookGetL2BadRequest()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	case 401:
		result := NewOrderBookGetL2Unauthorized()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	case 404:
		result := NewOrderBookGetL2NotFound()
		if err := result.readResponse(response, consumer, o.formats); err != nil {
			return nil, err
		}
		return nil, result

	default:
		return nil, runtime.NewAPIError("unknown error", response, response.Code())
	}
}

// NewOrderBookGetL2OK creates a OrderBookGetL2OK with default headers values
func NewOrderBookGetL2OK() *OrderBookGetL2OK {
	return &OrderBookGetL2OK{}
}

/*OrderBookGetL2OK handles this case with default header values.

Request was successful
*/
type OrderBookGetL2OK struct {
	Payload []*models.OrderBookL2
}

func (o *OrderBookGetL2OK) Error() string {
	return fmt.Sprintf("[GET /orderBook/L2][%d] orderBookGetL2OK  %+v", 200, o.Payload)
}

func (o *OrderBookGetL2OK) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	// response payload
	if err := consumer.Consume(response.Body(), &o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewOrderBookGetL2BadRequest creates a OrderBookGetL2BadRequest with default headers values
func NewOrderBookGetL2BadRequest() *OrderBookGetL2BadRequest {
	return &OrderBookGetL2BadRequest{}
}

/*OrderBookGetL2BadRequest handles this case with default header values.

Parameter Error
*/
type OrderBookGetL2BadRequest struct {
	Payload *models.Error
}

func (o *OrderBookGetL2BadRequest) Error() string {
	return fmt.Sprintf("[GET /orderBook/L2][%d] orderBookGetL2BadRequest  %+v", 400, o.Payload)
}

func (o *OrderBookGetL2BadRequest) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewOrderBookGetL2Unauthorized creates a OrderBookGetL2Unauthorized with default headers values
func NewOrderBookGetL2Unauthorized() *OrderBookGetL2Unauthorized {
	return &OrderBookGetL2Unauthorized{}
}

/*OrderBookGetL2Unauthorized handles this case with default header values.

Unauthorized
*/
type OrderBookGetL2Unauthorized struct {
	Payload *models.Error
}

func (o *OrderBookGetL2Unauthorized) Error() string {
	return fmt.Sprintf("[GET /orderBook/L2][%d] orderBookGetL2Unauthorized  %+v", 401, o.Payload)
}

func (o *OrderBookGetL2Unauthorized) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}

// NewOrderBookGetL2NotFound creates a OrderBookGetL2NotFound with default headers values
func NewOrderBookGetL2NotFound() *OrderBookGetL2NotFound {
	return &OrderBookGetL2NotFound{}
}

/*OrderBookGetL2NotFound handles this case with default header values.

Not Found
*/
type OrderBookGetL2NotFound struct {
	Payload *models.Error
}

func (o *OrderBookGetL2NotFound) Error() string {
	return fmt.Sprintf("[GET /orderBook/L2][%d] orderBookGetL2NotFound  %+v", 404, o.Payload)
}

func (o *OrderBookGetL2NotFound) readResponse(response runtime.ClientResponse, consumer runtime.Consumer, formats strfmt.Registry) error {

	o.Payload = new(models.Error)

	// response payload
	if err := consumer.Consume(response.Body(), o.Payload); err != nil && err != io.EOF {
		return err
	}

	return nil
}
